package qris

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/database/config"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

var qrisKey = "qris"

type Log struct {
	MpmRequest  *MpmRequest `json:"mpmRequest"`
	MpmResponse *MpmData    `json:"mpmResponse"`
	CpmRequest  *CpmRequest `json:"cpmRequest"`
	CpmResponse *CpmData    `json:"cpmResponse"`
}

func (qrisLog *Log) WriteToPostgres() {
	if db := config.InitMariaDB(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&qrisLog); err != nil {
			log.Error(err)
		}
		db.Create(&qrisLog)
	}
}

func (qrisLog *Log) WriteToMariaDB() {
	if db := config.InitPostgres(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&qrisLog); err != nil {
			log.Error(err)
		}
		db.Create(&qrisLog)
	}
}

func (qrisLog *Log) WriteToMongoDB() {
	if db, ctx, conf := config.InitMongoDB(); db != nil {
		defer config.CloseMongo(db, ctx)
		praySchedLog := db.Database(conf.Data).Collection(qrisKey)
		if bsonData, err := bson.Marshal(&praySchedLog); err != nil {
			log.Error(err)
		} else if _, err := praySchedLog.InsertOne(ctx, bsonData); err != nil {
			log.Error(err)
		}
	}
}

func (qrisLog *Log) WriteToLog() {
	qrisLog.WriteToMariaDB()
	qrisLog.WriteToPostgres()
	qrisLog.WriteToMongoDB()
}

func (qrisLog *Log) PutToRedisQueue() {
	config.PutToRedisQueue(qrisLog, qrisKey)
}

func (qrisLog *Log) GetFromRedis() {
	for {
		if rdb := config.InitRedis(); rdb != nil {
			ctx := context.Background()
			rawString := rdb.RPop(ctx, qrisKey).Val()
			if len(rawString) > 0 {
				qrisLog = &Log{}
				if err := json.Unmarshal([]byte(rawString), qrisLog); err != nil {
					log.Error(err)
				} else {
					qrisLog.WriteToLog()
				}
			}
			config.CloseRedis(rdb)
			sleepTime, _ := time.ParseDuration("1s")
			time.Sleep(sleepTime)
		}
	}
}
