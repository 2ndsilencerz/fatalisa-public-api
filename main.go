package main

import (
	"encoding/json"
	"fatalisa-public-api/config"
	"fatalisa-public-api/database"
	"github.com/pieterclaerhout/go-log"
)

func init() {
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006-01-02 15:04:05 -0700"
}

func init() {
	cfg := &config.List{}
	content, err := json.Marshal(cfg)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(content))

	cfgMariaDB := &database.MariaDBConf{}
	cfgMariaDB.Get()
	content, err = json.Marshal(cfgMariaDB)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(content))

	cfgPostgres := &database.PostgresConf{}
	cfgPostgres.Get()
	content, err = json.Marshal(cfgPostgres)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(content))
}

func init() {
	// run DB connection check async
	go database.DbConnCheck()
}

func main() {
	// any code after router.Run() won't be executed, place it in init() above
	log.Info("Starting service")
	router := &config.Router{}
	router.Run()
}
