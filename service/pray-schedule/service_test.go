package pray_schedule

import (
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/utils"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetSchedule(t *testing.T) {
	DownloadFile(83)
	req := &model.Request{
		City: "jakarta",
		Date: time.Now().Format("2006/01/02"),
	}
	res := getSchedule(req)
	tmp, _ := time.Parse("2006/01/02", req.Date)
	if res.Year != strconv.Itoa(tmp.Year()) ||
		res.Month != fmt.Sprintf("%02s", strconv.Itoa(int(tmp.Month()))) ||
		res.Date != fmt.Sprintf("%02s", strconv.Itoa(tmp.Day())) {
		t.Error("Data not match")
	}
}

func TestGetCityList(t *testing.T) {
	DownloadFile(83)
	jsonRes := utils.Jsonify(GetCityList())
	if !strings.Contains(jsonRes, "jakarta") {
		t.Error("Data not found")
	}
}