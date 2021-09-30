package config

import (
	"fatalisa-public-api/utils"
	"github.com/pieterclaerhout/go-log"
	"time"
)

func DbConnCheck() {
	for {
		log.Info("Checking DB connection")
		//go checkPostgres()
		//go checkMariaDB()
		go checkMongoDB()
		go checkRedis()
		sleepTime := utils.GetDuration("30s")
		time.Sleep(sleepTime)
	}
}

func printConf(v interface{}) {
	log.Debug(utils.Jsonify(v))
}
