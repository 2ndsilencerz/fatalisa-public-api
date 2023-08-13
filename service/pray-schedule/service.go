package pray_schedule

import (
	"encoding/xml"
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/subchen/go-log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var downloadTask = 0
var downloadGroup = sync.WaitGroup{}

const ScheduleFilesDir = "schedule/"
const filenamePrefix = "jadwal-"
const filenameExtension = ".xml"
const totalSchedules = 308
const maxSimultaneousDownloadTask = 3

var yearSchedule string

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

// GetSchedule used to get pray schedule of requested city and date if provided
func GetSchedule(c *fiber.Ctx) *model.Response {
	req := model.Request{}
	if err := c.BodyParser(&req); err != nil {
		req.City = c.Params("city")
		req.Date = time.Now().Format("2006/01/02")
	}
	log.Info(utils.Jsonify(req))

	res := getData(&req)
	log.Info(utils.Jsonify(res))
	return res
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
