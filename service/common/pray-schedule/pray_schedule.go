package pray_schedule

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fatalisa-public-api/database/config"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var HeaderPray = fmt.Sprintf("%-8s", "pray-sch")

func PraySchedDownload(duration string) {
	for {
		log.Info(HeaderPray, "|", "Downloading pray schedule")
		for i := 1; i <= 308; i++ {
			go DownloadFile(i)
		}
		sleepTime, err := time.ParseDuration(duration)
		if err != nil {
			log.Error(HeaderPray, "|", err)
		}
		time.Sleep(sleepTime)
	}
}

func DownloadFile(x int) {
	var data *Header
	cityCode := strconv.Itoa(x)

	url := "http://jadwalsholat.pkpu.or.id/export.php"
	contentType := "application/x-www-form-urlencoded"
	body := "period=3" + "&" +
		"y=2021" + "&" +
		"radio=1" + "&" +
		"fields_terminated=%3B" + "&" +
		"fields_enclosed=%22" + "&" +
		"lines_terminated=%5Cn%5Cr" + "&" +
		"edition=1" + "&" +
		"compress=0" + "&" +
		"adzanCountry=indonesia" + "&" +
		"adzanCity=" + cityCode + "&" +
		"language=indonesian" + "&" +
		"algo=1" + "&" +
		"cbxViewParam=1" + "&" +
		"cbxViewImsyak=1" + "&" +
		"cbxViewSunrise=1" + "&" +
		"observer_height=0" + "&" +
		"fajr=0" + "&" +
		"fajr_depr=20.0" + "&" +
		"fajr_interval=90.0" + "&" +
		"ashr=0" + "&" +
		"ashr_shadow=1.0" + "&" +
		"isha=0" + "&" +
		"isha_depr=18.0" + "&" +
		"isha_interval=90.0" + "&" +
		"imsyak_depr=1.5" + "&" +
		"imsyak=1" + "&" +
		"imsyak_interval=10.0" + "&" +
		"cmd=save" + "&" +
		"save=Print%2FCetak"
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Download
	res, err := client.Post(url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Error(HeaderPray, "|", err)
	}

	// Create file
	fileName := "jadwal-" + strconv.Itoa(x) + ".xml"
	file, errFileCreate := os.Create(fileName)
	if errFileCreate != nil {
		log.Error(HeaderPray, "|", errFileCreate)
	} else {
		body := res.Body
		if _, errWriteFile := io.Copy(file, body); errWriteFile != nil {
			log.Error(HeaderPray, "|", errWriteFile)
		}
		errCloseFile := file.Close()
		if errCloseFile != nil {
			log.Error(HeaderPray, "|", errCloseFile)
		}

		// Rename the file
		data = readFile(fileName)
		newFileName := "jadwal-" + data.City + ".xml"
		errRenameFile := os.Rename(fileName, newFileName)
		if errRenameFile != nil {
			log.Error(HeaderPray, "|", errRenameFile)
		}
		log.Info(HeaderPray, "|", "Pray schedule", data.City, "downloaded")
	}

	closeResForDownload(res)
}

func readFile(fileName string) *Header {
	res := &Header{}
	if file, errRead := os.ReadFile(fileName); errRead != nil {
		log.Error(HeaderPray, "|", errRead)
	} else if errParse := xml.Unmarshal(file, res); errParse != nil {
		log.Error(HeaderPray, "|", errParse)
	}
	return res
}

func closeResForDownload(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Error(HeaderPray, "|", err)
	}
}

type Header struct {
	Adzan `xml:"adzan"`
}

type Adzan struct {
	Version   string `xml:"version"`
	Site      string `xml:"site"`
	Country   string `xml:"country"`
	City      string `xml:"city"`
	Parameter `xml:"parameter"`
	Data      []PrayScheduleData `xml:"data" json:"data"`
}

type Parameter struct {
	Longitude string `xml:"longitude"`
	Latitude  string `xml:"latitude"`
	Direction string `xml:"direction"`
	Distance  string `xml:"distance"`
}

type PrayScheduleData struct {
	Year    string `xml:"year" json:"year" binding:"required"`
	Month   string `xml:"month" json:"month" binding:"required"`
	Date    string `xml:"date" json:"date" binding:"required"`
	Imsyak  string `xml:"imsyak" json:"imsyak" binding:"required"`
	Fajr    string `xml:"fajr" json:"fajr" binding:"required"`
	Syuruq  string `xml:"syuruq" json:"syuruq" binding:"required"`
	Dzuhur  string `xml:"dzuhur" json:"dzuhur" binding:"required"`
	Ashr    string `xml:"ashr" json:"ashr" binding:"required"`
	Maghrib string `xml:"maghrib" json:"maghrib" binding:"required"`
	Isha    string `xml:"isha" json:"isha" binding:"required"`
}

type PrayScheduleReq struct {
	City string `json:"city" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type PrayScheduleLog struct {
	gorm.Model
	UUID             uuid.UUID `json:"uuid" gorm:"column:uuid" bson:"uuid"`
	PrayScheduleReq  `json:"request" gorm:"column:request" bson:"request"`
	PrayScheduleData `json:"response" gorm:"column:response" bson:"response"`
}

type city struct {
	CityName string `json:"cityName"`
}

func GetSchedule(req *PrayScheduleReq) *PrayScheduleData {
	responseData := &PrayScheduleData{}
	fileName := "jadwal-" + req.City + ".xml"
	if data := readFile(fileName); data.Version != "" {
		if data.City == req.City {
			for i := 0; i < len(data.Data); i++ {
				date, _ := time.Parse("2006/01/02", req.Date)
				if data.Data[i].Year == date.Format("2006") &&
					data.Data[i].Month == date.Format("01") &&
					data.Data[i].Date == date.Format("02") {
					responseData = &data.Data[i]
				}
			}
		}
	}
	go saveLogToDB(*req, *responseData)
	return responseData
}

func GetCityList() interface{} {
	result := struct {
		List []*city `json:"list"`
	}{}
	files, err := os.ReadDir(".")
	if err != nil {
		log.Error(HeaderPray, "|", err)
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), "jadwal-") {
				log.Info(file.Name())
				cityName := strings.ReplaceAll(file.Name(), "jadwal-", "")
				cityName = strings.ReplaceAll(cityName, ".xml", "")
				result.List = append(result.List, &city{
					CityName: cityName,
				})
			}
		}
	}
	return result
}

func saveLogToDB(req PrayScheduleReq, res PrayScheduleData) {
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(HeaderPray, "|", err)
	}
	dbLog := &PrayScheduleLog{
		UUID:             uuidGenerated,
		PrayScheduleReq:  req,
		PrayScheduleData: res,
	}
	//dbLog.WriteToLog()
	config.PutToRedisQueue(dbLog, HeaderPray)
}

func (praySchedLog *PrayScheduleLog) WriteToPostgres() {
	if db := config.InitMariaDB(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(praySchedLog, "|", err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMariaDB() {
	if db := config.InitPostgres(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(praySchedLog, "|", err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMongoDB() {
	if db, ctx, conf := config.InitMongoDB(); db != nil {
		defer config.CloseMongo(db, ctx)
		accessLogCol := db.Database(conf.Data).Collection(HeaderPray)
		if bsonData, err := bson.Marshal(&praySchedLog); err != nil {
			log.Error(config.HeaderMongoDB, "|", err)
		} else {
			if _, err := accessLogCol.InsertOne(ctx, bsonData); err != nil {
				log.Error(config.HeaderMongoDB, "|", err)
			}
		}
	}
}

func (praySchedLog *PrayScheduleLog) WriteToLog() {
	praySchedLog.WriteToMariaDB()
	praySchedLog.WriteToPostgres()
	praySchedLog.WriteToMongoDB()
}

func (praySchedLog *PrayScheduleLog) GetFromRedis() {
	if rdb := config.InitRedis(); rdb != nil {
		defer config.CloseRedis(rdb)
		for {
			ctx := context.Background()
			rawString := rdb.RPop(ctx, HeaderPray).Val()
			if len(rawString) > 0 {
				praySchedLog = &PrayScheduleLog{}
				if err := json.Unmarshal([]byte(rawString), praySchedLog); err != nil {
					log.Error(HeaderPray, "|", err)
				} else {
					praySchedLog.WriteToLog()
				}
			}
			sleepTime, _ := time.ParseDuration("1s")
			time.Sleep(sleepTime)
		}
	}
}
