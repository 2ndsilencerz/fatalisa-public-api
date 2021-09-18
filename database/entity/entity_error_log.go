package entity

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/database/config"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"time"
)

type ErrorLog struct {
	UUID      uuid.UUID `json:"uuid" bson:"uuid"`
	Message   string    `json:"message" bson:"message"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

var ErrorLogKey = "error_log"

func (errorLog *ErrorLog) WriteToMariaDB() {
	if db := config.InitMariaDB(); db != nil {
		if err := db.AutoMigrate(&errorLog); err != nil {
			log.Error(config.HeaderGorm, "|", err)
		}
		db.Create(&errorLog)
		config.CloseGorm(db)
	}
}

func (errorLog *ErrorLog) WriteToPostgres() {
	if db := config.InitPostgres(); db != nil {
		if err := db.AutoMigrate(&errorLog); err != nil {
			log.Error(config.HeaderGorm, "|", err)
		}
		db.Create(&errorLog)
		config.CloseGorm(db)
	}
}

func (errorLog *ErrorLog) Write(err error) {
	errorLog.Message = err.Error()
	errorLog.Timestamp = time.Now()
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(config.HeaderGorm, "|", err)
	}
	errorLog.UUID = uuidGenerated
	config.PutToRedisQueue(&errorLog, ErrorLogKey)
}

func (errorLog *ErrorLog) WriteLog() {
	errorLog.WriteToMariaDB()
	errorLog.WriteToPostgres()
}

func (errorLog *ErrorLog) GetFromRedis() {
	for {
		if rdb := config.InitRedis(); rdb != nil {
			ctx := context.Background()
			rawString := rdb.RPop(ctx, AccessLogKey).Val()
			if len(rawString) > 0 {
				errorLog = &ErrorLog{}
				if err := json.Unmarshal([]byte(rawString), errorLog); err != nil {
					log.Error(ErrorLogKey, "|", err)
				} else {
					errorLog.WriteLog()
				}
			}
			config.CloseRedis(rdb)
			sleepTime, _ := time.ParseDuration("1s")
			time.Sleep(sleepTime)
		}
	}
}
