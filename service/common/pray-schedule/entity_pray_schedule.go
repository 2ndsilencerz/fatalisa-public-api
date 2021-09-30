package pray_schedule

import (
	"fatalisa-public-api/database/config"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

var praySchedKey = "pray_schedule_log"

type Tabler interface {
	TableName() string
}

func (PrayScheduleLog) TableName() string {
	return praySchedKey
}

type PrayScheduleLog struct {
	gorm.Model
	UUID             uuid.UUID `json:"uuid" gorm:"column:uuid" bson:"uuid"`
	PrayScheduleReq  `json:"request" gorm:"column:request" bson:"request"`
	PrayScheduleData `json:"response" gorm:"column:response" bson:"response"`
	Created          time.Time `json:"created"`
}

func (praySchedLog *PrayScheduleLog) WriteToPostgres() {
	db := config.InitMariaDB()
	db.Write(praySchedLog)
}

func (praySchedLog *PrayScheduleLog) WriteToMariaDB() {
	db := config.InitPostgres()
	db.Write(praySchedLog)
}

func (praySchedLog *PrayScheduleLog) WriteToMongoDB() {
	db := config.InitMongoDB()
	db.InsertOne(praySchedKey, praySchedLog)
}

func (praySchedLog *PrayScheduleLog) WriteToLog() {
	praySchedLog.WriteToMariaDB()
	praySchedLog.WriteToPostgres()
	praySchedLog.WriteToMongoDB()
}

func (praySchedLog *PrayScheduleLog) PutToRedisQueue() {
	rdb := config.InitRedis()
	rdb.PushQueue(praySchedKey, praySchedLog)
}

func (praySchedLog *PrayScheduleLog) GetFromRedis() {
	for {
		rdb := config.InitRedis()
		if rdb.PopQueue(praySchedKey, praySchedLog); len(praySchedLog.Year) > 0 {
			praySchedLog.WriteToLog()
		}
		// since they use same address for storing the data, we need to reinstate
		// so the next data fetched will be fresh
		praySchedLog = &PrayScheduleLog{}
		//sleepTime := utils.GetDuration("1s")
		//time.Sleep(sleepTime)
	}
}
