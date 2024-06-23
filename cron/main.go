package main

import (
	"context"
	"fatalisa-public-api/database/config"
	"fatalisa-public-api/service/pray-schedule/jadwalsholatorg"
	"fatalisa-public-api/service/pray-schedule/model"
	"github.com/subchen/go-log"
	"time"
)

// This is a onetime run to download all pray schedule data
// use this with cron
func main() {
	//praySchedulePkpuSvc.PrayScheduleDownloadPKPU()
	ctx := context.Background()
	cityList := jadwalsholatorg.GetCityList(ctx)
	today := time.Now().Format("2006/01/02")
	redis := config.InitRedis()
	for _, city := range cityList.List {
		req := model.Request{}
		req.City = city.Name
		req.Date = today
		if deleted := redis.Delete("schedule:"+req.City, ctx); !deleted {
			log.Warn("Existing cache failed to delete or non-existent")
		}
		jadwalsholatorg.GetSchedule(&req, ctx)
	}
}
