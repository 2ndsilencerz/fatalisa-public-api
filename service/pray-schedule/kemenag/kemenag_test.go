package kemenag

import (
	"fmt"
	"github.com/subchen/go-log"
	"strconv"
	"testing"
	"time"
)

const (
	cityTestKemenag     = "KOTA JAKARTA"
	provinceTestKemenag = "DKI JAKARTA"
)

var err error

func TestGetProvince(t *testing.T) {
	res := *GetProvinces()
	log.Info(res)
	if string(res[provinceTestKemenag]) == "" {
		t.Error()
		return
	}
}

func TestGetCity(t *testing.T) {
	resProvinces := *GetProvinces()
	log.Info(resProvinces)
	resCities := *GetCities(string(resProvinces[provinceTestKemenag]))
	log.Info(resCities)
	if string(resCities[cityTestKemenag]) == "" {
		t.Error()
		return
	}
}

func TestGetSchedule(t *testing.T) {
	resProvinces := *GetProvinces()
	log.Info(resProvinces[provinceTestKemenag])
	resCities := *GetCities(string(resProvinces[provinceTestKemenag]))
	if resCities[cityTestKemenag] == "" {
		t.Error()
		return
	}
	var thisMonth, thisYear int
	thisMonth, err = strconv.Atoi(time.Now().Format("01"))
	if err != nil {
		t.Error(err)
	}
	thisYear, err = strconv.Atoi(time.Now().Format("2006"))
	if err != nil {
		t.Error(err)
	}
	log.Info(provinceTestKemenag + " " + cityTestKemenag + " " + string(resCities[cityTestKemenag]) + " " +
		fmt.Sprint(thisMonth) + " " + fmt.Sprint(thisYear))
	resSchedule := GetSchedule(string(resProvinces[provinceTestKemenag]), string(resCities[cityTestKemenag]),
		thisMonth, thisYear)
	if resSchedule.Subuh == "" {
		t.Error()
		return
	}
}
