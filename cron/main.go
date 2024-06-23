package main

import (
	"context"
	"fatalisa-public-api/service/pray-schedule/jadwalsholatorg"
	"fatalisa-public-api/service/pray-schedule/model"
	"time"
)

// This is a onetime run to download all pray schedule data
// use this with cron
func main() {
	//praySchedulePkpuSvc.PrayScheduleDownloadPKPU()
	ctx := context.Background()
	cityList := jadwalsholatorg.GetCityList(ctx)
	today := time.Now().Format("2006/01/02")
	for _, city := range cityList.List {
		req := model.Request{}
		req.City = city.Name
		req.Date = today
		jadwalsholatorg.GetSchedule(&req, ctx)
	}
}
