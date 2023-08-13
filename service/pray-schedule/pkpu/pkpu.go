package pkpu

import (
	"bytes"
	"encoding/xml"
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const totalSchedules = 308
const maxSimultaneousDownloadTask = 3
const ScheduleFilesDir = "schedule/"
const filenamePrefix = "jadwal-"
const filenameExtension = ".xml"

var downloadTask = 0
var downloadGroup = sync.WaitGroup{}
var yearSchedule string

// Header is the header of xml used to parse schedule data downloaded
type Header struct {
	Adzan `xml:"adzan"`
}

// Adzan used as model to fetch schedule data from xml
type Adzan struct {
	Version   string `xml:"version"`
	Site      string `xml:"site"`
	Country   string `xml:"country"`
	City      string `xml:"city"`
	Parameter `xml:"parameter"`
	Data      []model.Response `xml:"data" json:"data"`
}

// Parameter used to get misc data from schedule in xml
type Parameter struct {
	Longitude string `xml:"longitude"`
	Latitude  string `xml:"latitude"`
	Direction string `xml:"direction"`
	Distance  string `xml:"distance"`
}

// PrayScheduleDownload used in cronjob for downloading all cities schedule
func PrayScheduleDownload() {
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

func GetDataPkpu(cityCode int) *http.Response {
	cityCodeStr := strconv.Itoa(cityCode)

	url := "http://jadwalsholat.pkpu.or.id/export.php"
	contentType := "application/code-www-form-urlencoded"
	body := "period=3" + "&" + // 3 for all year schedule
		"y=" + yearSchedule + "&" + // year selection
		"radio=1" + "&" +
		"fields_terminated=%3B" + "&" +
		"fields_enclosed=%22" + "&" +
		"lines_terminated=%5Cn%5Cr" + "&" +
		"edition=1" + "&" +
		"compress=0" + "&" +
		"adzanCountry=indonesia" + "&" +
		"adzanCity=" + cityCodeStr + "&" + // city selection
		"language=indonesian" + "&" + // id language selection
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
		log.Error(err)
		return nil
	}
	return res
}

// DownloadFile used to download the pray schedule based on the city code
/*
City code used here are an int starting with 1,
and it's mapping isn't available since the source are from http://jadwalsholat.pkpu.or.id/
*/
func DownloadFile(cityCode int) {
	var data *Header

	// Download
	if res := GetDataPkpu(cityCode); res != nil {
		// Create dir
		utils.CheckAndCreateDir(ScheduleFilesDir)

		// Create file
		fileName := utils.GetWorkingDir() + ScheduleFilesDir + filenamePrefix + strconv.Itoa(cityCode) + filenameExtension
		file := utils.CreateFile(fileName)
		if file != nil {
			body := res.Body
			if _, errWriteFile := io.Copy(file, body); errWriteFile != nil {
				log.Error("Write to file error: ", errWriteFile)
			}
			//if errOwnFile := file.Chown(1001, 1001); errOwnFile != nil {
			//	log.Error("Change owner of file error: ", errOwnFile)
			//}
			//if errModFile := file.Chmod(664); errModFile != nil {
			//	log.Error("Change access mode of file error: ", errModFile)
			//}
			if errCloseFile := file.Close(); errCloseFile != nil {
				log.Error(errCloseFile)
			}

			// Rename the file postfix from city cityCode to actual city name
			data = readFile(fileName)
			newFileName := ScheduleFilesDir + filenamePrefix + data.City + filenameExtension
			log.Info("Renaming file from ", fileName, " to ", filenamePrefix+data.City+filenameExtension)
			if errRenameFile := os.Rename(fileName, newFileName); errRenameFile != nil {
				log.Error("Rename file error: ", errRenameFile)
			}
			log.Info("Pray schedule ", data.City, " downloaded")
		}
		closeResponseForDownload(res)
	}

	if downloadTask > 0 {
		downloadGroup.Done()
		downloadTask--
	}
}

func GetData(req *model.Request) *model.Response {
	responseData := model.Response{}
	fileName := ScheduleFilesDir + filenamePrefix + req.City + filenameExtension
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

func readFile(fileName string) *Header {
	res := Header{}
	if file, errRead := os.ReadFile(fileName); errRead != nil {
		log.Error(errRead)
	} else if errParse := xml.Unmarshal(file, &res); errParse != nil {
		log.Error(errParse)
	}
	return &res
}

func closeResponseForDownload(response *http.Response) {
	if err := response.Body.Close(); err != nil {
		log.Error(err)
	}
}

// GetCityList used to get city list that can be used to get pray schedule
func GetCityList() interface{} {
	res := model.CityList{}
	if files, err := os.ReadDir(ScheduleFilesDir); err != nil {
		log.Error(err)
	} else {
		for _, file := range files {
			if strings.Contains(file.Name(), filenamePrefix) {
				cityName := strings.ReplaceAll(file.Name(), filenamePrefix, "")
				cityName = strings.ReplaceAll(cityName, filenameExtension, "")
				city := &model.City{}
				city.Name = cityName
				res.List = append(res.List, city)
			}
		}
	}
	log.Info(utils.Jsonify(res))
	return &res
}
