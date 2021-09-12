package database

import (
	"time"
)

type ErrorLog struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func (errorLog ErrorLog) WriteToMariaDB() {
	if db := InitMariaDB(); db != nil {
		db.Create(&errorLog)
		Close(db)
	}
}

func (errorLog ErrorLog) WriteToPostgres() {
	if db := InitPostgres(); db != nil {
		db.Create(&errorLog)
		Close(db)
	}
}

func (errorLog ErrorLog) Write(err error) {
	errorLog.Message = err.Error()
	errorLog.Timestamp = time.Now()
	errorLog.WriteToMariaDB()
	errorLog.WriteToPostgres()
}
