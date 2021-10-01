package pray_schedule

import (
	"bytes"
	"encoding/xml"
	"fatalisa-public-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var downloadTask = 0
var downloadGroup = sync.WaitGroup{}
var scheduleHourDownload = []int{0, 6, 12, 18}

const filenamePrefix = "schedule/jadwal-"
const filenameExtension = ".xml"
const totalSchedules = 308

const maxSimultaneousDownloadTask = 3

var yearSchedule string

func checkTime(currentTime time.Time) bool {
	for _, v := range scheduleHourDownload {
		if v == currentTime.Hour() && currentTime.Minute() == 0 && currentTime.Second() == 0 {
			return true
		}
	}
	return false
}

func scheduleFilesExist() bool {
	if fileList, err := ioutil.ReadDir("/schedule"); err != nil {
		log.Error(err)
	} else if len(fileList) > 0 && len(fileList) == totalSchedules-1 {
		return true
	}
	return false
}

func PraySchedDownload() {
	for {
		currentTime := time.Now()
		yearSchedule = strconv.Itoa(currentTime.Year())
		if checkTime(currentTime) || !scheduleFilesExist() {
			log.Info("Downloading pray schedule")
			for i := 1; i <= totalSchedules; i++ {
				downloadTask++
				downloadGroup.Add(1)
				go DownloadFile(i)
				if downloadTask > maxSimultaneousDownloadTask {
					downloadGroup.Wait()
				}
			}
		}
	}
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
		fileName := filenamePrefix + strconv.Itoa(x) + filenameExtension
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
			newFileName := filenamePrefix + data.City + filenameExtension
			if errRenameFile := os.Rename(fileName, newFileName); errRenameFile != nil {
				log.Error(errRenameFile)
			}
			log.Info("Pray schedule", data.City, "downloaded")
		}
		closeResForDownload(res)
	}

	if downloadTask > 0 {
		downloadGroup.Done()
		downloadTask--
	}
}

func readFile(fileName string) *Header {
	res := &Header{}
	if file, errRead := os.ReadFile(fileName); errRead != nil {
		log.Error(errRead)
	} else if errParse := xml.Unmarshal(file, res); errParse != nil {
		log.Error(errParse)
	}
	return res
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
	responseData := &PrayScheduleData{}
	fileName := filenamePrefix + req.City + filenameExtension
	if data := readFile(fileName); data.Version != "" && data.City == req.City {
		for i := 0; i < len(data.Data); i++ {
			date, _ := time.Parse("2006/01/02", req.Date)
			if data.Data[i].Year == date.Format("2006") &&
				data.Data[i].Month == date.Format("01") &&
				data.Data[i].Date == date.Format("02") {
				responseData = &data.Data[i]
			}
		}
	}
	return responseData
}

func GetScheduleService(c *gin.Context) *PrayScheduleData {
	req := &PrayScheduleReq{}
	if err := c.BindJSON(req); err != nil {
		log.Error(err)
	} else {
		log.Info(utils.Jsonify(req))
	}

	res := getSchedule(req)
	log.Info(utils.Jsonify(res))
	return res
}

func GetCityList() interface{} {
	res := &CityListRes{}
	if files, err := os.ReadDir("."); err != nil {
		log.Error(err)
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), "jadwal-") {
				cityName := strings.ReplaceAll(file.Name(), "jadwal-", "")
				cityName = strings.ReplaceAll(cityName, ".xml", "")
				city := &city{}
				city.CityName = cityName
				res.List = append(res.List, city)
			}
		}
	}
	log.Info(utils.Jsonify(res))
	return res
}
