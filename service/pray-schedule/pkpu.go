package pray_schedule

import (
	"bytes"
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

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

func GetDataPkpu(x int) *http.Response {
	cityCode := strconv.Itoa(x)

	url := "http://jadwalsholat.pkpu.or.id/export.php"
	contentType := "application/x-www-form-urlencoded"
	body := "period=3" + "&" + // 3 for all year schedule
		"y=" + yearSchedule + "&" + // year selection
		"radio=1" + "&" +
		"fields_terminated=%3B" + "&" +
		"fields_enclosed=%22" + "&" +
		"lines_terminated=%5Cn%5Cr" + "&" +
		"edition=1" + "&" +
		"compress=0" + "&" +
		"adzanCountry=indonesia" + "&" +
		"adzanCity=" + cityCode + "&" + // city selection
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
func DownloadFile(x int) {
	var data *Header

	// Download
	if res := GetDataPkpu(x); res != nil {
		// Create dir
		if _, err := os.Stat(ScheduleFilesDir); err != nil {
			log.Warn(err)
			workDir := utils.GetWorkingDir()
			createDir := workDir + ScheduleFilesDir
			utils.Mkdir(createDir)
		}

		// Create file
		fileName := utils.GetWorkingDir() + ScheduleFilesDir + filenamePrefix + strconv.Itoa(x) + filenameExtension
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

			// Rename the file postfix from city code to actual city name
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

func getData(req *model.Request) *model.Response {
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
