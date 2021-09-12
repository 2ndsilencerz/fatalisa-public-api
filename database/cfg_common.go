package database

import (
	"fmt"
	"github.com/pieterclaerhout/go-log"
	"gorm.io/gorm"
	"time"
)

var HeaderGorm = fmt.Sprintf("%-8s", "gorm")

func Close(database *gorm.DB) {
	if db, err := database.DB(); err != nil {
		log.Error(HeaderGorm, "|", err)
		//panic(err)
	} else {
		if err = db.Close(); err != nil {
			log.Error(HeaderGorm, "|", err)
			//panic(err)
		}
	}
}

func DbConnCheck() {
	for {
		log.Info(HeaderGorm, "|", "Checking DB connection")
		go checkPostgres()
		go checkMariaDB()
		if sleepTime, err := time.ParseDuration("30s"); err != nil {
			log.Error(HeaderGorm, "|", err)
		} else {
			time.Sleep(sleepTime)
		}
	}
}

func checkPostgres() {
	if postgres := InitPostgres(); postgres != nil {
		Close(postgres)
	}
}

func checkMariaDB() {
	if mariadb := InitMariaDB(); mariadb != nil {
		Close(mariadb)
	}
}
