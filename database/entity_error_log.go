package database

import (
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"time"
)

type ErrorLog struct {
	UUID      uuid.UUID `json:"uuid" bson:"uuid"`
	Message   string    `json:"message" bson:"message"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

func (errorLog *ErrorLog) WriteToMariaDB() {
	if db := InitMariaDB(); db != nil {
		if err := db.AutoMigrate(&errorLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&errorLog)
		Close(db)
	}
}

func (errorLog *ErrorLog) WriteToPostgres() {
	if db := InitPostgres(); db != nil {
		if err := db.AutoMigrate(&errorLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&errorLog)
		Close(db)
	}
}

func (errorLog *ErrorLog) Write(err error) {
	errorLog.Message = err.Error()
	errorLog.Timestamp = time.Now()
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(HeaderGorm, "|", err)
	}
	errorLog.UUID = uuidGenerated
	errorLog.WriteToMariaDB()
	errorLog.WriteToPostgres()
}
