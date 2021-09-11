package database

import (
	"github.com/pieterclaerhout/go-log"
	"gorm.io/gorm"
	"time"
)

func Close(database *gorm.DB) {
	db, err := database.DB()
	if err != nil {
		log.Error(err)
		//panic(err)
	}
	err = db.Close()
	if err != nil {
		log.Error(err)
		//panic(err)
	}
}

func DbConnCheck() {
	for {
		log.Info("Checking DB connection")
		go checkPostgres()
		go checkMariaDB()
		sleepTime, err := time.ParseDuration("30s")
		if err != nil {
			log.Error(err)
			errorLog := &ErrorLog{}
			errorLog.Write(err)
		}
		time.Sleep(sleepTime)
	}
}

func checkPostgres() {
	postgres := InitPostgres()
	Close(postgres)
}

func checkMariaDB() {
	mariadb := InitMariaDB()
	Close(mariadb)
}
