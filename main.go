package main

import (
	"encoding/json"
	"fatalisa-public-api/config"
	"fatalisa-public-api/config/router"
	"fatalisa-public-api/database"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/pieterclaerhout/go-log"
)

func init() {
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006/01/02 15:04:05 -0700"
}

var HeaderJSON = fmt.Sprintf("%-8s", "json")
var headerCfg = fmt.Sprintf("%-8s", "svc-cfg")

func init() {
	// run DB connection check async
	go database.DbConnCheck()
	// run scheduled download for certain time
	go utils.ScheduleDownload("1h")
}

//func init() {
//	for k, v := range os.Environ() {
//		log.Info(k, v)
//	}
//}

func init() {
	cfg := &config.List{}
	cfg.Get()
	if content, err := json.Marshal(cfg); err != nil {
		log.Error(HeaderJSON, "|", err)
	} else {
		log.Info(headerCfg, "|", string(content))
	}
}

func main() {
	// any code after routerInit.Run() won't be executed, place it in init() above
	log.Info(router.HeaderGin, "|", "Starting service")
	routerInit := &router.Router{}
	routerInit.Run()
}
