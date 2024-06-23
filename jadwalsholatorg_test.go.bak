package main

import (
	"context"
	"fatalisa-public-api/service/pray-schedule/jadwalsholatorg"
	"fatalisa-public-api/service/pray-schedule/model"
	"github.com/subchen/go-log"
	"strings"
	"testing"
	"time"
)

func TestGetCityList(t *testing.T) {
	resCities := jadwalsholatorg.GetCityList(context.Background())
	log.Info(resCities.List[0].Name)

	if len(resCities.List) == 0 {
		t.Fail()
	}
}

func TestGetSchedule(t *testing.T) {
	req := model.Request{}
	city := "Aceh Barat"
	req.City = city
	req.Date = time.Now().Format("2006/01/02")
	resSchedule := jadwalsholatorg.GetSchedule(&req, context.Background())

	if resSchedule.Fajr == "" || !strings.EqualFold(resSchedule.City, city) {
		t.Fail()
	}
}
