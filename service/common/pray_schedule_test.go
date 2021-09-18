package common

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestGetSchedule(t *testing.T) {
	DownloadFile()
	req := &PrayScheduleReq{
		City: "jakarta",
		Date: time.Now().Format("2006/01/02"),
	}
	res := GetSchedule(req)
	tmp, _ := time.Parse("2006/01/02", req.Date)
	if res.Year != strconv.Itoa(tmp.Year()) ||
		res.Month != fmt.Sprintf("%02s", strconv.Itoa(int(tmp.Month()))) ||
		res.Date != fmt.Sprintf("%02s", strconv.Itoa(tmp.Day())) {
		t.Error("Data not match")
	}
}
