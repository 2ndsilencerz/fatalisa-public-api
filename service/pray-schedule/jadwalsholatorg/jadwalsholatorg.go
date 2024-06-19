package jadwalsholatorg

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/database/config"
	"fatalisa-public-api/service/pray-schedule/model"
	utils2 "fatalisa-public-api/service/web/utils"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/subchen/go-log"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "https://jadwalsholat.org"
const monthlyUrl = baseUrl + "/jadwal-sholat/monthly.php"
const dateUrl = "?id="
const monthUrl = "&m="
const yearUrl = "&y="

const cityRedisKey = "city:"
const scheduleRedisKey = "schedule:"

var timeList = []string{"subuh", "zuhur", "ashar", "maghrib", "isya"}

var redis *config.RedisConf

func init() {
	redis = config.InitRedis()
}

func GetCityList(ctx context.Context) model.CityList {
	res := model.CityList{}
	if citiesRedis := redis.GetKeys(cityRedisKey+"*", ctx); len(citiesRedis) > 0 {
		for key := range citiesRedis {
			city := &model.City{}
			city.Name = strings.Replace(key, cityRedisKey, "", -1)
			res.List = append(res.List, city)
		}
		return res
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resSchedule, err := client.Get(monthlyUrl)
	if err, _ := utils2.ErrorHandler(err); err {
		return res
	}

	doc, err := html.Parse(resSchedule.Body)
	if err, _ := utils2.ErrorHandler(err); err {
		return res
	}
	var extractData func(*html.Node)
	extractData = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "option" {
			for _, a := range n.Attr {
				cityCode := a.Val
				city := &model.City{}
				city.Name = n.FirstChild.Data

				//log.Info(a.Val + " " + n.FirstChild.Data)
				redis.Set(cityRedisKey+city.Name, cityCode, ctx, 0)
				res.List = append(res.List, city)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractData(c)
		}
	}
	extractData(doc)

	log.Info(utils.Jsonify(res))
	return res
}

func GetSchedule(req *model.Request, ctx context.Context) *model.Response {
	res := model.Response{}
	today := time.Now().Format("2006-01-02")
	dateReq, err := time.Parse("2006-01-02", req.Date)
	if err, _ := utils2.ErrorHandler(err); err {
		return &res
	}

	if len(redis.GetKeys(cityRedisKey+"*", ctx)) == 0 {
		GetCityList(ctx)
	}

	if strings.EqualFold(req.Date, today) {
		log.Info("Checking cache for today's schedule")
		//log.Info(scheduleRedisKey + req.City)
		scheduleCityToday := redis.Get(scheduleRedisKey+req.City, ctx)
		if scheduleCityToday != nil {
			scheduleMap := make(map[string]string)
			if err, _ := utils2.ErrorHandler(json.Unmarshal([]byte(fmt.Sprint(scheduleCityToday)), &scheduleMap)); !err {
				log.Info("Schedule found in cache")
				res.Year = strconv.Itoa(dateReq.Year())
				res.Month = strconv.Itoa(int(dateReq.Month()))
				res.Date = strconv.Itoa(dateReq.Day())

				res.Syuruq = fmt.Sprint(scheduleMap[timeList[0]])
				res.Dzuhur = fmt.Sprint(scheduleMap[timeList[1]])
				res.Ashr = fmt.Sprint(scheduleMap[timeList[2]])
				res.Maghrib = fmt.Sprint(scheduleMap[timeList[3]])
				res.Isha = fmt.Sprint(scheduleMap[timeList[4]])
				return &res
			}
		}
		log.Warn("Failed to get today's schedule from cache")
	}

	cityCode := fmt.Sprint(redis.Get(cityRedisKey+req.City, ctx))
	if len(cityCode) == 0 {
		log.Warn("City Code not found in cache, will default search to Jakarta Pusat")
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	url := monthlyUrl + dateUrl + cityCode + monthUrl + strconv.Itoa(int(dateReq.Month())) + yearUrl + strconv.Itoa(dateReq.Year())
	log.Info(url)
	resSchedule, err := client.Get(url)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}

	scheduleList := make(map[string]string)
	doc, err := html.Parse(resSchedule.Body)
	if err, _ := utils2.ErrorHandler(err); err {
		return nil
	}
	log.Info(doc)

	var extractData func(*html.Node)
	extractData = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			// <tr><td><b>
			currentNode := n.FirstChild
			dateNode := n.FirstChild.FirstChild
			// only get for specific date
			if dateNode != nil && dateNode.Data == "b" && dateNode.FirstChild != nil && fmt.Sprint(dateNode.FirstChild.Data) == fmt.Sprintf("%02d", dateReq.Day()) {
				//Imsyak -> Shubuh -> Terbit -> Dhuha -> Dzuhur -> Ashr -> Maghrib -> Isya
				subuhNode := currentNode.NextSibling.NextSibling
				zuhurNode := subuhNode.NextSibling.NextSibling.NextSibling
				asharNode := zuhurNode.NextSibling
				maghribNode := asharNode.NextSibling
				isyaNode := maghribNode.NextSibling
				scheduleList[timeList[0]] = subuhNode.FirstChild.Data
				scheduleList[timeList[1]] = zuhurNode.FirstChild.Data
				scheduleList[timeList[2]] = asharNode.FirstChild.Data
				scheduleList[timeList[3]] = maghribNode.FirstChild.Data
				scheduleList[timeList[4]] = isyaNode.FirstChild.Data
				//log.Info(fmt.Sprint(scheduleList))
				jsonMap := utils.Jsonify(scheduleList)
				//if err, _ := utils2.ErrorHandler(err); !err {
				redis.Set(scheduleRedisKey+req.City, jsonMap, ctx, 0)
				//}

				res.Year = strconv.Itoa(dateReq.Year())
				res.Month = strconv.Itoa(int(dateReq.Month()))
				res.Date = strconv.Itoa(dateReq.Day())

				res.Syuruq = scheduleList[timeList[0]]
				res.Dzuhur = scheduleList[timeList[1]]
				res.Ashr = scheduleList[timeList[2]]
				res.Maghrib = scheduleList[timeList[3]]
				res.Isha = scheduleList[timeList[4]]
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractData(c)
		}
	}
	extractData(doc)

	return &res
}
