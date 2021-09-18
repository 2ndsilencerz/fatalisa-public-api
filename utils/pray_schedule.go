package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fatalisa-public-api/database"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var HeaderPray = fmt.Sprintf("%-8s", "pray-sch")

func ScheduleDownload(duration string, waitGroup *sync.WaitGroup) {
	for {
		waitGroup.Add(1)
		log.Info(HeaderPray, "|", "Downloading pray schedule")
		downloadFile()
		waitGroup.Done()
		sleepTime, err := time.ParseDuration(duration)
		if err != nil {
			log.Error(HeaderPray, "|", err)
		}
		time.Sleep(sleepTime)
	}
}

func downloadFile() {
	fileName := "jadwal.xml"
	file, errFileCreate := os.Create(fileName)
	if errFileCreate != nil {
		log.Error(HeaderPray, "|", errFileCreate)
	}
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
		"adzanCity=83" + "&" +
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
	// Put content on file
	res, err := client.Post(url, contentType, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.Error(HeaderPray, "|", err)
	}
	//res, errDownload := http.Post(url, contentType, bytes.NewBuffer([]byte(body)))
	//if errDownload != nil {
	//	log.Error(errDownload)
	//}
	defer closeResForDownload(res)
	if errFileCreate == nil {
		body := res.Body
		//os.WriteFile(fileName, , 0777)
		if _, errWriteFile := io.Copy(file, body); errWriteFile != nil {
			log.Error(HeaderPray, "|", errWriteFile)
		} else {
			log.Info(HeaderPray, "|", "Pray schedule downloaded")
		}
	}
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

func GetSchedule(req *PrayScheduleReq) *PrayScheduleData {
	responseData := &PrayScheduleData{}
	if file, errRead := ioutil.ReadFile("jadwal.xml"); errRead != nil {
		log.Error(HeaderPray, "|", errRead)
	} else {
		//data := &headerXML{}
		data := &Header{}
		if errParse := xml.Unmarshal(file, data); errParse != nil {
			log.Error(HeaderPray, "|", errParse)
		} else {
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
	}
	go saveLogToDB(*req, *responseData)
	return responseData
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
	database.PutToRedisQueue(dbLog, HeaderPray)
}

func (praySchedLog *PrayScheduleLog) WriteToPostgres() {
	if db := database.InitMariaDB(); db != nil {
		defer database.Close(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(praySchedLog, "|", err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMariaDB() {
	if db := database.InitPostgres(); db != nil {
		defer database.Close(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(praySchedLog, "|", err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMongoDB() {
	if db, ctx, conf := database.InitMongoDB(); db != nil {
		defer database.CloseMongo(db, ctx)
		accessLogCol := db.Database(conf.Data).Collection(HeaderPray)
		if bsonData, err := bson.Marshal(&praySchedLog); err != nil {
			log.Error(database.HeaderMongoDB, "|", err)
		} else {
			if _, err := accessLogCol.InsertOne(ctx, bsonData); err != nil {
				log.Error(database.HeaderMongoDB, "|", err)
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
	if rdb := database.InitRedis(); rdb != nil {
		defer database.CloseRedis(rdb)
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
