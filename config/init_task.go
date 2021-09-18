package config

import (
	dbCfg "fatalisa-public-api/database/config"
	"fatalisa-public-api/database/entity"
	svc "fatalisa-public-api/service/common"
	"github.com/pieterclaerhout/go-log"
	"os"
)

//var headerCfg = fmt.Sprintf("%-8s", "svc-cfg")

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
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006/01/02 15:04:05 -0700"
}

// run scheduled DB check on another routine
func init() {
	// run DB connection check async
	go dbCfg.DbConnCheck()
}

// run redis queue checker per entity
func init() {
	accessLog := &entity.AccessLog{}
	go accessLog.GetFromRedis()

	praySchedLog := &svc.PrayScheduleLog{}
	go praySchedLog.GetFromRedis()

	errorLog := &entity.ErrorLog{}
	go errorLog.GetFromRedis()
}

// run scheduled download for certain time
func init() {
	go svc.PraySchedDownload("1h")
}

//// print config list
//func init() {
//	cfg := &config.List{}
//	cfg.Get()
//	if content, err := json.Marshal(cfg); err != nil {
//		log.Error(headerCfg, "|", err)
//	} else {
//		log.Info(headerCfg, "|", string(content))
//	}
//}
