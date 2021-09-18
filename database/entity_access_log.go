package database

import (
	"context"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"gorm.io/gorm"
	"time"
)

/*
Column name convention used in Gorm.io framework is snake_case
In-case you want to assign with different name, use tag `gorm:"column:columnName"`
*/

type AccessLog struct {
	gorm.Model
	UUID       uuid.UUID `json:"uuid" gorm:"column:uuid" bson:"uuid"`
	Kind       string    `json:"message" gorm:"column:kind" bson:"kind"`
	IP         string    `json:"ip" bson:"ip"`
	Method     string    `json:"method" bson:"method"`
	FullPath   string    `json:"full_path" bson:"full_path"`
	StatusCode int       `json:"status_code" bson:"status_code"`
	Created    int64     `gorm:"autoCreateTime,column:created" json:"created" bson:"created"`
}

var AccessLogKey = "access_log"

func (accessLog *AccessLog) WriteToMariaDB() {
	if db := InitMariaDB(); db != nil {
		defer Close(db)
		if err := db.AutoMigrate(&accessLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&accessLog)
	}
}

func (accessLog *AccessLog) WriteToPostgres() {
	if db := InitPostgres(); db != nil {
		defer Close(db)
		if err := db.AutoMigrate(&accessLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&accessLog)
	}
}

func (accessLog *AccessLog) WriteToMongoDB() {
	if db, ctx, conf := InitMongoDB(); db != nil {
		defer CloseMongo(db, ctx)
		accessLogCol := db.Database(conf.Data).Collection(AccessLogKey)
		if bsonData, err := bson.Marshal(&accessLog); err != nil {
			log.Error(HeaderMongoDB, "|", err)
		} else {
			if _, err := accessLogCol.InsertOne(ctx, bsonData); err != nil {
				log.Error(HeaderMongoDB, "|", err)
			}
		}
	}
}

func (accessLog *AccessLog) WriteLog() {
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(HeaderGorm, "|", err)
	}
	accessLog.UUID = uuidGenerated
	accessLog.WriteToMariaDB()
	accessLog.WriteToPostgres()
	accessLog.WriteToMongoDB()
}

func (accessLog *AccessLog) GetFromRedis() {
	if rdb := InitRedis(); rdb != nil {
		defer CloseRedis(rdb)
		for {
			ctx := context.Background()
			rawString := rdb.RPop(ctx, AccessLogKey).Val()
			if len(rawString) > 0 {
				accessLog = &AccessLog{}
				if err := json.Unmarshal([]byte(rawString), accessLog); err != nil {
					log.Error(AccessLogKey, "|", err)
				} else {
					accessLog.WriteLog()
				}
			}
			sleepTime, _ := time.ParseDuration("1s")
			time.Sleep(sleepTime)
		}
	}
}

//func (accessLog *AccessLog) InitSubscriber() *redis.PubSub {
//	ctx := context.Background()
//	rdb := InitRedis()
//	subscriber := rdb.Subscribe(ctx, AccessLogKey)
//	return subscriber
//}
//
//func subscribe(subscriber *redis.PubSub) {
//	//for {
//		ctx := context.Background()
//		iface, err := subscriber.Receive(ctx)
//		if err != nil {
//			log.Error(HeaderRedis, "|", err)
//		}
//		switch msg := iface.(type) {
//		case *redis.Subscription:
//			log.Info(HeaderRedis, "|", "Subscribed to channel", msg.Channel)
//			break
//		case *redis.Message:
//			accessLog := &AccessLog{}
//			if err := json.Unmarshal([]byte(msg.Payload), accessLog); err != nil {
//				log.Error(HeaderRedis, "|", err)
//			}
//			accessLog.WriteLog()
//			break
//		//case *redis.Pong:
//		//	break
//		default:
//			log.Error(HeaderRedis, "|", "Something happened")
//		}
//		_ = subscriber.Channel()
//	//}
//}
