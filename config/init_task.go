package config

import (
	prayScheduleSvc "fatalisa-public-api/service/common/pray-schedule"
	"github.com/pieterclaerhout/go-log"
	"os"
)

func Init() {
}

// print BUILD_DATE if exist
func init() {
	if buildDate, exist := os.LookupEnv("BUILD_DATE"); exist && len(buildDate) > 0 {
		log.Info("This image built in", buildDate)
	}
}

// set logger format
func init() {
	log.DebugMode = false
	log.DebugSQLMode = false
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006/01/02 15:04:05 -0700"
}

// run scheduled DB check on another routine
//func init() {
//	// run DB connection check async
//	go dbCfg.DbConnCheck()
//}

// run redis queue checker per entity
//func init() {
//	accessLog := &router.AccessLog{}
//	go accessLog.GetFromRedis()
//
//	praySchedLog := &prayScheduleSvc.PrayScheduleLog{}
//	go praySchedLog.GetFromRedis()
//
//	//errorLog := &common.ErrorLog{}
//	//go errorLog.GetFromRedis()
//
//	qrisLog := &qrisSvc.Log{}
//	go qrisLog.GetFromRedis()
//}

// run scheduled download for certain time
func init() {
	go prayScheduleSvc.PraySchedDownload()
}
