package kemenag

import (
	"encoding/json"
	internalModel "fatalisa-public-api/service/pray-schedule/kemenag/model"
	"fatalisa-public-api/service/pray-schedule/model"
	utils2 "fatalisa-public-api/service/web/utils"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/subchen/go-log"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	baseUrl     = "https://bimasislam.kemenag.go.id"
	getProvince = baseUrl + "/jadwalshalat"
	getCity     = baseUrl + "/ajax/getKabkoshalat"
	getSchedule = baseUrl + "/ajax/getShalatbln"

	ScheduleFilesDir = "schedule/"
	provinceLocation = ScheduleFilesDir + "province.json"
	cityLocation     = ScheduleFilesDir + "city-{}.json"
)

var (
	Phpsessid         string
	BimasislamSession string
)

func Init() {
	provinces := GetProvinces()
	if _, err := os.Stat(utils.GetWorkingDir() + provinceLocation); err != nil {
		go saveProvinceMap(provinces)
	}
	for _, provinceCode := range provinces {
		cities := GetCities(string(provinceCode))
		if _, err := os.Stat(utils.GetWorkingDir() + strings.Replace(cityLocation, "{}", string(provinceCode), -1)); err != nil {
			go saveCityMap(cities, string(provinceCode))
		}
	}
	// for initializing Phpsessid and BimasislamSession values
	go requestProvinces(nil)
}

func requestProvinces(response *http.Response) *http.Response {
	var (
		res *http.Response
		err error
	)
	urlProvince := getProvince
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	req, err := http.NewRequest(http.MethodGet, urlProvince, nil)
	if response != nil && Phpsessid == "" && BimasislamSession == "" {
		cookies := response.Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "PHPSESSID" {
				Phpsessid = cookie.Value
			}
			if cookie.Name == "bimasislam_session" {
				BimasislamSession = cookie.Value
			}
		}
		// set cookies
		req = setCookies(req)
		// execute request
		res, err = client.Do(req)
	} else {
		res, err = client.Do(req)
		// do request once more because the first response isn't complete
		res = requestProvinces(res)
	}
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	return res
}

func requestCities(provinceCode string) *http.Response {
	urlCities := getCity
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	req, err := http.NewRequest(http.MethodPost, urlCities, strings.NewReader("x="+provinceCode))
	req = setCookies(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	return res
}

func requestSchedule(provinceCode, cityCode string, month, year int) *http.Response {
	strMonth := strconv.Itoa(month)
	strYear := strconv.Itoa(year)
	urlSchedule := getSchedule
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	req, err := http.NewRequest(http.MethodPost, urlSchedule,
		strings.NewReader("x="+provinceCode+"&y="+cityCode+"&bln="+strMonth+"&thn="+strYear))
	req = setCookies(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	return res
}

func setCookies(req *http.Request) *http.Request {
	// set cookies
	req.AddCookie(&http.Cookie{
		Name:  "PHPSESSID",
		Value: Phpsessid,
	})
	req.AddCookie(&http.Cookie{
		Name:  "bimasislam_session",
		Value: BimasislamSession,
	})
	//req.Header.SetString("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	return req
}

func GetProvinces() internalModel.ProvinceMapping {
	var province internalModel.ProvinceMapping
	if province = readProvinceMap(); province == nil {
		response := requestProvinces(nil)
		htmlContent, _ := html.Parse(response.Body)
		province = make(internalModel.ProvinceMapping)
		province.ParseHtml(htmlContent)
	}
	return province
}

func GetCities(provinceCode string) internalModel.CityMapping {
	var cities internalModel.CityMapping
	if cities = readCityMap(provinceCode); cities == nil {
		response := requestCities(provinceCode)
		htmlContent, _ := html.Parse(response.Body)
		cities = make(internalModel.CityMapping)
		cities.ParseHtml(htmlContent)
	}
	return cities
}

func GetSchedule(provinceCode, cityCode string, month, year int) internalModel.Daily {
	schedule := GetScheduleMonth(provinceCode, cityCode, month, year)
	today := fmt.Sprint(year) + "-" + fmt.Sprintf("%02d", month) + "-" + time.Now().Format("02")
	log.Info(today)
	return schedule.Data[today]
}

func GetScheduleMonth(provinceCode, cityCode string, month, year int) *internalModel.Schedule {
	response := requestSchedule(provinceCode, cityCode, month, year)
	schedule := new(internalModel.Schedule)
	schedule.Parse(response)
	return schedule
}

func GetData(req *model.Request) *model.Response {
	responseData := model.Response{}
	var (
		province, city string
		month, year    int
	)
	// fix city request
	req.City, _ = url.QueryUnescape(req.City)
	req.City = strings.TrimSpace(req.City)
	req.City = strings.ToUpper(req.City)
	//log.Info(req.City)

	provinces := GetProvinces()
	for _, provinceCode := range provinces {

		cities := GetCities(string(provinceCode))
		if cities[internalModel.CityName(req.City)] != "" {
			province = string(provinceCode)
			city = string(cities[internalModel.CityName(req.City)])
			break
		}
	}
	dateRequested, _ := time.Parse("2006/01/02", req.Date)
	month, _ = strconv.Atoi(dateRequested.Format("01"))
	year, _ = strconv.Atoi(dateRequested.Format("2006"))
	//theDate := fmt.Sprint(year) + "-" + fmt.Sprintf("%02d", month) + "-" + fmt.Sprintf("%02d", dateRequested.Day())
	//log.Info(theDate)
	dailyData := GetSchedule(province, city, month, year)
	responseData.Year = fmt.Sprint(year)
	responseData.Month = fmt.Sprintf("%02d", month)
	responseData.Date = fmt.Sprintf("%02d", dateRequested.Day())
	responseData.Syuruq = dailyData.Subuh
	responseData.Dzuhur = dailyData.Dzuhur
	responseData.Ashr = dailyData.Ashar
	responseData.Maghrib = dailyData.Maghrib
	responseData.Isha = dailyData.Isya
	return &responseData
}

// GetCityList used to get city list that can be used to get pray schedule
func GetCityList() interface{} {
	res := model.CityList{}
	provinces := GetProvinces()
	for _, provinceCode := range provinces {
		cities := GetCities(string(provinceCode))

		for city := range cities {
			currCity := &model.City{}
			currCity.Name = string(city)
			res.List = append(res.List, currCity)
		}
	}
	log.Info(utils.Jsonify(res))
	return &res
}

func saveProvinceMap(mapping internalModel.ProvinceMapping) {
	jsonMapping, err := json.Marshal(mapping)
	if err, _ := utils2.ErrorHandler(err); err {
		return
	}
	utils.CheckAndCreateDir(ScheduleFilesDir)
	utils2.ErrorHandler(os.WriteFile(utils.GetWorkingDir()+provinceLocation, jsonMapping, 0644))
}

func saveCityMap(mapping internalModel.CityMapping, provinceCode string) {
	jsonMapping, err := json.Marshal(mapping)
	if err, _ := utils2.ErrorHandler(err); err {
		return
	}
	utils.CheckAndCreateDir(ScheduleFilesDir)
	fileName := utils.GetWorkingDir() + strings.Replace(cityLocation, "{}", provinceCode, -1)
	utils2.ErrorHandler(os.WriteFile(fileName, jsonMapping, 0644))
}

func readProvinceMap() internalModel.ProvinceMapping {
	file, err := os.ReadFile(utils.GetWorkingDir() + provinceLocation)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	var mapping internalModel.ProvinceMapping
	err = json.Unmarshal(file, &mapping)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	return mapping
}

func readCityMap(provinceCode string) internalModel.CityMapping {
	file, err := os.ReadFile(utils.GetWorkingDir() + strings.Replace(cityLocation, "{}", provinceCode, -1))
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	var mapping internalModel.CityMapping
	err = json.Unmarshal(file, &mapping)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	return mapping
}
