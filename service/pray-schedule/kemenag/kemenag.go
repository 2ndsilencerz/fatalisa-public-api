package kemenag

import (
	"encoding/json"
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/subchen/go-log"
	"golang.org/x/net/html"
	"io"
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

type (
	ProvinceCode    string
	ProvinceName    string
	ProvinceMapping map[ProvinceName]ProvinceCode

	CityCode    string
	CityName    string
	CityMapping map[CityName]CityCode

	Schedule struct {
		Status  int     `json:"status"`
		Message string  `json:"message"`
		Prov    string  `json:"prov"`
		Kabko   string  `json:"kabko"`
		Data    Dailies `json:"data"`
	}

	Dailies map[string]Daily

	Daily struct {
		Tanggal string `json:"tanggal"`
		Imsak   string `json:"imsak"`
		Subuh   string `json:"subuh"`
		Terbit  string `json:"terbit"`
		Dhuha   string `json:"dhuha"`
		Dzuhur  string `json:"dzuhur"`
		Ashar   string `json:"ashar"`
		Maghrib string `json:"maghrib"`
		Isya    string `json:"isya"`
	}
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
	if err != nil {
		log.Error(err)
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
	if err != nil {
		log.Error(err)
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
	if err != nil {
		log.Error(err)
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
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36")
	return req
}

func (province ProvinceMapping) parseHtml(node *html.Node) {
	if node.Parent != nil && node.Parent.Data == "select" && len(node.Parent.Attr) > 0 &&
		node.Parent.Attr[0].Val == "search_prov" && node.Data == "option" {
		if !strings.HasPrefix(node.FirstChild.Data, "PILIH") {
			province[ProvinceName(node.FirstChild.Data)] = ProvinceCode(node.Attr[0].Val)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		province.parseHtml(c)
	}
}

func (city CityMapping) parseHtml(node *html.Node) {
	if node.Data == "option" {
		city[CityName(node.FirstChild.Data)] = CityCode(node.Attr[0].Val)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		city.parseHtml(c)
	}
}

func (schedule *Schedule) parse(response *http.Response) {
	// parse response to Schedule
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal(raw, &schedule)
	if err != nil {
		log.Error(err)
		return
	}
}

func GetProvinces() ProvinceMapping {
	var province ProvinceMapping
	if province = readProvinceMap(); province == nil {
		response := requestProvinces(nil)
		htmlContent, _ := html.Parse(response.Body)
		province = make(ProvinceMapping)
		province.parseHtml(htmlContent)
	}
	return province
}

func GetCities(provinceCode string) CityMapping {
	var cities CityMapping
	if cities = readCityMap(provinceCode); cities == nil {
		response := requestCities(provinceCode)
		htmlContent, _ := html.Parse(response.Body)
		cities = make(CityMapping)
		cities.parseHtml(htmlContent)
	}
	return cities
}

func GetSchedule(provinceCode, cityCode string, month, year int) Daily {
	schedule := GetScheduleMonth(provinceCode, cityCode, month, year)
	today := fmt.Sprint(year) + "-" + fmt.Sprintf("%02d", month) + "-" + time.Now().Format("02")
	log.Info(today)
	return schedule.Data[today]
}

func GetScheduleMonth(provinceCode, cityCode string, month, year int) *Schedule {
	response := requestSchedule(provinceCode, cityCode, month, year)
	schedule := new(Schedule)
	schedule.parse(response)
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
	//log.Info(req.City)

	provinces := GetProvinces()
	for _, provinceCode := range provinces {

		cities := GetCities(string(provinceCode))
		if cities[CityName(req.City)] != "" {
			province = string(provinceCode)
			city = string(cities[CityName(req.City)])
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

func saveProvinceMap(mapping ProvinceMapping) {
	jsonMapping, err := json.Marshal(mapping)
	if err != nil {
		log.Error(err)
		return
	}
	utils.CheckAndCreateDir(ScheduleFilesDir)
	err = os.WriteFile(utils.GetWorkingDir()+provinceLocation, jsonMapping, 0644)
	if err != nil {
		log.Error(err)
	}
}

func saveCityMap(mapping CityMapping, provinceCode string) {
	jsonMapping, err := json.Marshal(mapping)
	if err != nil {
		log.Error(err)
		return
	}
	utils.CheckAndCreateDir(ScheduleFilesDir)
	fileName := utils.GetWorkingDir() + strings.Replace(cityLocation, "{}", provinceCode, -1)
	err = os.WriteFile(fileName, jsonMapping, 0644)
	if err != nil {
		log.Error(err)
	}
}

func readProvinceMap() ProvinceMapping {
	file, err := os.ReadFile(utils.GetWorkingDir() + provinceLocation)
	if err != nil {
		log.Error(err)
		return nil
	}
	var mapping ProvinceMapping
	err = json.Unmarshal(file, &mapping)
	if err != nil {
		log.Error(err)
		return nil
	}
	return mapping
}

func readCityMap(provinceCode string) CityMapping {
	file, err := os.ReadFile(utils.GetWorkingDir() + strings.Replace(cityLocation, "{}", provinceCode, -1))
	if err != nil {
		log.Error(err)
		return nil
	}
	var mapping CityMapping
	err = json.Unmarshal(file, &mapping)
	if err != nil {
		log.Error(err)
		return nil
	}
	return mapping
}
