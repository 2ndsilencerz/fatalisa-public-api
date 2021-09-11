package main

import (
	"encoding/json"
	"fatalisa-public-api/config"
	"fatalisa-public-api/database"
	"github.com/pieterclaerhout/go-log"
)

func main() {
	cfg := &config.List{}
	cfg.Get()
	content, err := json.Marshal(cfg)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(content))

	// run DB connection check async
	go database.DbConnCheck()
}

func init() {
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006-01-02 15:04:05 -0700"

	log.Info("Starting service")

	router := &config.Router{}
	router.Run()
}
