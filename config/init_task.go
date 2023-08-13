package config

import (
	"fatalisa-public-api/service/common"
	prayschedule "fatalisa-public-api/service/pray-schedule/pkpu"
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
func init() {
	logDir := utils.GetWorkingDir() + utils.FileLogLocation
	if _, err := os.Stat(logDir); err != nil {
		log.Warn(err)
		utils.Mkdir(logDir)
	}
	scheduleDir := utils.GetWorkingDir() + prayschedule.ScheduleFilesDir
	if _, err := os.Stat(scheduleDir); err != nil {
		log.Warn(err)
		utils.Mkdir(scheduleDir)
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
