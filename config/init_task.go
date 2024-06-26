package config

import (
	"fatalisa-public-api/service/common"
	pray_schedule "fatalisa-public-api/service/pray-schedule"
	"fatalisa-public-api/service/pray-schedule/kemenag"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"os"
)

func Init() {
}

// set logger format
func init() {
	//log.Default.Out = io.MultiWriter(
	//	os.Stdout,
	//	&writers.DailyFileWriter{
	//		Name:     utils.FileLogName,
	//		MaxCount: 10,
	//	},
	//)
	log.Default.Formatter = new(utils.LogFormat)
}

// set timezone by looking the environment variable
func init() {
	_, _ = os.LookupEnv("TZ")
}

// print BUILD_DATE if exist
func init() {
	if txt := common.VersionChecker().Message; len(txt) > 0 {
		log.Info("This image built in ", txt)
	}
}

// check directory for service if existed
// when it's not, create one
//func init() {
//	logDir := utils.FileLogLocation
//	utils.CheckAndCreateDir(logDir)
//	scheduleDir := prayschedule.ScheduleFilesDir
//	utils.CheckAndCreateDir(scheduleDir)
//}

func init() {
	provider, _ := os.LookupEnv("PROVIDER")
	if provider == pray_schedule.ProviderKemenag {
		kemenag.Init()
	}
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
