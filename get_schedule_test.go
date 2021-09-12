package main

import (
	"bytes"
	"encoding/json"
	"fatalisa-public-api/config/router"
	"fatalisa-public-api/utils"
	"github.com/pieterclaerhout/go-log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouter() *router.Router {
	routerInit := &router.Router{}
	routerInit.Get()
	routerInit.InitRoutes()
	return routerInit
}

func TestGetSchedule(t *testing.T) {
	routerTest := setupRouter()

	current := time.Now()
	city := "jakarta"
	bodyReq := &utils.PrayScheduleReq{
		City: city,
		Date: current.Format("2006/01/02"),
	}

	if bodyReqJson, err := json.Marshal(bodyReq); err != nil {
		t.Error(err)
	} else {
		httpRes := httptest.NewRecorder()
		httpReq, err := http.NewRequest("POST", "/api/pray-schedule", bytes.NewBuffer(bodyReqJson))
		if err != nil {
			log.Error(err)
		} else {
			routerTest.R.ServeHTTP(httpRes, httpReq)

			dataRes := &utils.PrayScheduleData{}
			if rawRes, err := ioutil.ReadAll(httpRes.Body); err != nil {
				log.Error(err)
			} else {
				if err := json.Unmarshal(rawRes, dataRes); err != nil {
					log.Error(err)
				}
				if current.Format("2006") != dataRes.Year &&
					current.Format("01") != dataRes.Month &&
					current.Format("02") != dataRes.Date {
					t.Errorf("Wrong date, expected %s got %s/%s/%s",
						current.Format("2006/01/02"), dataRes.Year, dataRes.Month, dataRes.Date)
				}
			}
		}
	}
}
