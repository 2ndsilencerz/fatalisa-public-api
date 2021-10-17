package config

import (
	pray_schedule "fatalisa-public-api/service/common/pray-schedule"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"github.com/subchen/go-log/writers"
	"io"
	"os"
)

func Init() {
}

func init() {
	_, _ = os.LookupEnv("TZ")
}

// print BUILD_DATE if exist
func init() {
	if buildDate, exist := os.LookupEnv("BUILD_DATE"); exist && len(buildDate) > 0 {
		log.Info("This image built in ", buildDate)
	}
}

// check directory if exist
func init() {
	if _, err := os.Stat(utils.FileLogLocation); err != nil {
		log.Error(err)
		utils.Mkdir(utils.FileLogLocation)
	}
	if _, err := os.Stat(pray_schedule.ScheduleFilesDir); err != nil {
		log.Error(err)
		utils.Mkdir(pray_schedule.ScheduleFilesDir)
	}
}

// set logger format
func init() {
	log.Default.Out = io.MultiWriter(
		os.Stdout,
		&writers.DailyFileWriter{
			Name:     utils.FileLogName,
			MaxCount: 10,
		},
	)
	log.Default.Formatter = new(utils.LogFormat)
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

// run scheduled download
//func init() {
//	go prayScheduleSvc.ScheduleDownload()
//}
