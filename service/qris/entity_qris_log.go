package qris

import (
	"fatalisa-public-api/database/config"
	"fatalisa-public-api/utils"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

var qrisKey = "qris"

type Tabler interface {
	TableName() string
}

func (Log) TableName() string {
	return qrisKey
}

type Log struct {
	gorm.Model
	UUID        uuid.UUID   `json:"uuid" bson:"uuid"`
	MpmRequest  *MpmRequest `json:"mpmRequest" bson:"mpmRequest" gorm:"embedded"`
	MpmResponse *MpmData    `json:"mpmResponse" bson:"mpmResponse" gorm:"embedded"`
	CpmRequest  *CpmRequest `json:"cpmRequest" bson:"cpmRequest" gorm:"embedded"`
	CpmResponse *CpmData    `json:"cpmResponse" bson:"cpmResponse" gorm:"embedded"`
	Created     time.Time   `json:"created"`
}

func (qrisLog *Log) WriteToPostgres() {
	db := config.InitMariaDB()
	db.Write(qrisLog)
}

func (qrisLog *Log) WriteToMariaDB() {
	db := config.InitPostgres()
	db.Write(qrisLog)
}

func (qrisLog *Log) WriteToMongoDB() {
	db := config.InitMongoDB()
	db.InsertOne(qrisKey, qrisLog)
}

func (qrisLog *Log) WriteToLog() {
	qrisLog.WriteToMariaDB()
	qrisLog.WriteToPostgres()
	qrisLog.WriteToMongoDB()
}

func (qrisLog *Log) PutToRedisQueue() {
	rdb := config.InitRedis()
	rdb.PushQueue(qrisKey, qrisLog)
}

func (qrisLog *Log) GetFromRedis() {
	for {
		rdb := config.InitRedis()
		if rdb.PopQueue(qrisKey, qrisLog); len(utils.Jsonify(qrisLog.MpmRequest)) > 2 || len(utils.Jsonify(qrisLog.CpmRequest)) > 2 {
			qrisLog.WriteToLog()
		}
		// since they use same address for storing the data, we need to reinstate
		// so the next data fetched will be fresh
		qrisLog = &Log{}
		//sleepTime := utils.GetDuration("1s")
		//time.Sleep(sleepTime)
	}
}
