package common

import (
	"fatalisa-public-api/database/config"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"time"
)

var errorLogKey = "error_log"

type Tabler interface {
	TableName() string
}

func (ErrorLog) TableName() string {
	return errorLogKey
}

type ErrorLog struct {
	UUID      uuid.UUID `json:"uuid" bson:"uuid"`
	Message   string    `json:"message" bson:"message"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

var ErrorLogKey = "error_log"

func (errorLog *ErrorLog) WriteToMariaDB() {
	db := config.InitMariaDB()
	db.Write(errorLog)
}

func (errorLog *ErrorLog) WriteToPostgres() {
	db := config.InitPostgres()
	db.Write(errorLog)
}

func (errorLog *ErrorLog) PutToRedisQueue(err error) {
	errorLog.Message = err.Error()
	errorLog.Timestamp = time.Now()
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(err)
	}
	errorLog.UUID = uuidGenerated
	rdb := config.InitRedis()
	rdb.PushQueue(ErrorLogKey, &errorLog)
}

func (errorLog *ErrorLog) WriteLog() {
	errorLog.WriteToMariaDB()
	errorLog.WriteToPostgres()
}

func (errorLog *ErrorLog) GetFromRedis() {
	for {
		rdb := config.InitRedis()
		if rdb.PopQueue(errorLogKey, errorLog); len(errorLog.Message) > 0 {
			errorLog.WriteLog()
		}
		// since they use same address for storing the data, we need to reinstate
		// so the next data fetched will be fresh
		errorLog = &ErrorLog{}
		//sleepTime := utils.GetDuration("1s")
		//time.Sleep(sleepTime)
	}
}
