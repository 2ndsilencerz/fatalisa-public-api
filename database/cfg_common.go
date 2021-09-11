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
		panic(err)
	}
	err = db.Close()
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func DbConnCheck() {
	for {
		postgres := InitPostgres()
		Close(postgres)
		mariadb := InitMariaDB()
		Close(mariadb)
		sleepTime, err := time.ParseDuration("10s")
		if err != nil {
			log.Error(err)
			errorLog := &ErrorLog{}
			errorLog.Write(err)
		}
		time.Sleep(sleepTime)
	}
}
