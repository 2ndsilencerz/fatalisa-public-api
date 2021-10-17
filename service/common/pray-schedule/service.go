package pray_schedule

import (
	"bytes"
	"encoding/xml"
	"fatalisa-public-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var downloadTask = 0
var downloadGroup = sync.WaitGroup{}

const scheduleFilesDir = "schedule/"
const filenamePrefix = "jadwal-"
const filenameExtension = ".xml"
const totalSchedules = 308

const maxSimultaneousDownloadTask = 3

var yearSchedule string

//func ScheduleDownload() {
//	var err error
//	s1 := gocron.NewScheduler(time.UTC)
//	s2 := gocron.NewScheduler(time.UTC)
//	s3 := gocron.NewScheduler(time.UTC)
//	s4 := gocron.NewScheduler(time.UTC)
//	time1, _ := time.Parse("15:04:05", "00:00:00")
//	time2, _ := time.Parse("15:04:05", "06:00:00")
//	time3, _ := time.Parse("15:04:05", "12:00:00")
//	time4, _ := time.Parse("15:04:05", "18:00:00")
//	if _, err = s1.Every(1).Days().StartAt(time1).Do(PraySchedDownload); err != nil {
//		log.Error(err)
//	}
//	if _, err = s2.Every(1).Days().StartAt(time2).Do(PraySchedDownload); err != nil {
//		log.Error(err)
//	}
//	if _, err = s3.Every(1).Days().StartAt(time3).Do(PraySchedDownload); err != nil {
//		log.Error(err)
//	}
//	if _, err = s4.Every(1).Days().StartAt(time4).Do(PraySchedDownload); err != nil {
//		log.Error(err)
//	}
//	s1.StartAsync()
//	s2.StartAsync()
//	s3.StartAsync()
//	s4.StartAsync()
//}

func PraySchedDownload() {
	log.Infof("%s %v", "Downloading pray schedule at", time.Now())
	for i := 1; i <= totalSchedules; i++ {
		downloadTask++
		downloadGroup.Add(1)
		go DownloadFile(i)
		if downloadTask > maxSimultaneousDownloadTask {
			downloadGroup.Wait()
		}
	}
	downloadGroup.Wait()
}

func DownloadFile(x int) {
	var data *Header
	cityCode := strconv.Itoa(x)

	url := "http://jadwalsholat.pkpu.or.id/export.php"
	contentType := "application/x-www-form-urlencoded"
	body := "period=3" + "&" +
		"y=" + yearSchedule + "&" +
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
	if res, err := client.Post(url, contentType, bytes.NewBuffer([]byte(body))); err != nil {
		log.Error(err)
	} else {
		// Create file
		fileName := scheduleFilesDir + filenamePrefix + strconv.Itoa(x) + filenameExtension
		file, errFileCreate := os.Create(fileName)
		if errFileCreate != nil {
			log.Error(errFileCreate)
		} else {
			body := res.Body
			if _, errWriteFile := io.Copy(file, body); errWriteFile != nil {
				log.Error(errWriteFile)
			}
			if errOwnFile := file.Chown(0, 0); errOwnFile != nil {
				log.Error(errOwnFile)
			}
			if errModFile := file.Chmod(0777); errModFile != nil {
				log.Error(errModFile)
			}
			if errCloseFile := file.Close(); errCloseFile != nil {
				log.Error(errCloseFile)
			}

			// Rename the file
			data = readFile(fileName)
			newFileName := scheduleFilesDir + filenamePrefix + data.City + filenameExtension
			if errRenameFile := os.Rename(fileName, newFileName); errRenameFile != nil {
				log.Error(errRenameFile)
			}
			log.Info("Pray schedule ", data.City, " downloaded")
		}
		closeResForDownload(res)
	}

	if downloadTask > 0 {
		downloadGroup.Done()
		downloadTask--
	}
}

func readFile(fileName string) *Header {
	res := Header{}
	if file, errRead := os.ReadFile(fileName); errRead != nil {
		log.Error(errRead)
	} else if errParse := xml.Unmarshal(file, &res); errParse != nil {
		log.Error(errParse)
	}
	return &res
}

func closeResForDownload(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Error(err)
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

type CityListRes struct {
	List []*city `json:"list"`
}

type city struct {
	CityName string `json:"cityName"`
}

func getSchedule(req *PrayScheduleReq) *PrayScheduleData {
	responseData := PrayScheduleData{}
	fileName := scheduleFilesDir + filenamePrefix + req.City + filenameExtension
	if data := readFile(fileName); data.Version != "" && data.City == req.City {
		for i := 0; i < len(data.Data); i++ {
			date, _ := time.Parse("2006/01/02", req.Date)
			if data.Data[i].Year == date.Format("2006") &&
				data.Data[i].Month == date.Format("01") &&
				data.Data[i].Date == date.Format("02") {
				responseData = data.Data[i]
			}
		}
	}
	return &responseData
}

func GetScheduleService(c *gin.Context) *PrayScheduleData {
	req := PrayScheduleReq{}
	// replace from BindJSON to ShouldBinJSON, so we should handle the error ourselves
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		log.Warn("Request method is GET")
		req.City = c.Param("city")
		req.Date = time.Now().Format("2006/01/02")
	}
	log.Info(utils.Jsonify(req))

	res := getSchedule(&req)
	log.Info(utils.Jsonify(res))
	return res
}

func GetCityList() interface{} {
	res := CityListRes{}
	if files, err := os.ReadDir(scheduleFilesDir); err != nil {
		log.Error(err)
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), filenamePrefix) {
				cityName := strings.ReplaceAll(file.Name(), filenamePrefix, "")
				cityName = strings.ReplaceAll(cityName, filenameExtension, "")
				city := &city{}
				city.CityName = cityName
				res.List = append(res.List, city)
			}
		}
	}
	log.Info(utils.Jsonify(res))
	return &res
}
