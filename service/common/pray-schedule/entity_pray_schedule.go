package pray_schedule

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/database/config"
	"fatalisa-public-api/utils"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/bson"
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
}

func (praySchedLog *PrayScheduleLog) WriteToPostgres() {
	if db := config.InitMariaDB(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMariaDB() {
	if db := config.InitPostgres(); db != nil {
		defer config.CloseGorm(db)
		if err := db.AutoMigrate(&praySchedLog); err != nil {
			log.Error(err)
		}
		db.Create(&praySchedLog)
	}
}

func (praySchedLog *PrayScheduleLog) WriteToMongoDB() {
	if db, ctx, conf := config.InitMongoDB(); db != nil {
		defer config.CloseMongo(db, ctx)
		praySchedLog := db.Database(conf.Data).Collection(praySchedKey)
		if bsonData, err := bson.Marshal(&praySchedLog); err != nil {
			log.Error(err)
		} else if _, err := praySchedLog.InsertOne(ctx, bsonData); err != nil {
			log.Error(err)
		}
	}
}

func (praySchedLog *PrayScheduleLog) WriteToLog() {
	praySchedLog.WriteToMariaDB()
	praySchedLog.WriteToPostgres()
	praySchedLog.WriteToMongoDB()
}

func (praySchedLog *PrayScheduleLog) PutToRedisQueue() {
	config.PutToRedisQueue(praySchedLog, praySchedKey)
}

func (praySchedLog *PrayScheduleLog) GetFromRedis() {
	for {
		if rdb := config.InitRedis(); rdb != nil {
			ctx := context.Background()
			rawString := rdb.RPop(ctx, praySchedKey).Val()
			if len(rawString) > 0 {
				praySchedLog = &PrayScheduleLog{}
				if err := json.Unmarshal([]byte(rawString), praySchedLog); err != nil {
					log.Error(err)
				} else {
					praySchedLog.WriteToLog()
				}
			}
			config.CloseRedis(rdb)
			sleepTime := utils.GetDuration("1s")
			time.Sleep(sleepTime)
		}
	}
}
