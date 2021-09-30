package router

import (
	"fatalisa-public-api/database/config"
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"gorm.io/gorm"
	"time"
)

var accessLogKey = "access_log"

type Tabler interface {
	TableName() string
}

func (AccessLog) TableName() string {
	return accessLogKey
}

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
	Created    time.Time `gorm:"column:created" json:"created" bson:"created"` //gorm:"autoCreateTime,column:created"
}

func (accessLog *AccessLog) WriteToMariaDB() {
	db := config.InitMariaDB()
	db.Write(accessLog)
}

func (accessLog *AccessLog) WriteToPostgres() {
	db := config.InitPostgres()
	db.Write(accessLog)
}

func (accessLog *AccessLog) WriteToMongoDB() {
	db := config.InitMongoDB()
	db.InsertOne(accessLogKey, accessLog)
}

func (accessLog *AccessLog) WriteLog() {
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(err)
	}
	accessLog.UUID = uuidGenerated
	accessLog.WriteToMariaDB()
	accessLog.WriteToPostgres()
	accessLog.WriteToMongoDB()
}

func (accessLog *AccessLog) PutToRedisQueue() {
	rdb := config.InitRedis()
	rdb.PushQueue(accessLogKey, accessLog)
}

func (accessLog *AccessLog) GetFromRedis() {
	for {
		rdb := config.InitRedis()
		rdb.PopQueue(accessLogKey, accessLog)
		if len(accessLog.Kind) > 0 {
			accessLog.WriteLog()
		}
		// since they use same address for storing the data, we need to reinstate
		// so the next data fetched will be fresh
		accessLog = &AccessLog{}
		//sleepTime := utils.GetDuration("1s")
		//time.Sleep(sleepTime)
	}
}

//func (accessLog *AccessLog) InitSubscriber() *redis.PubSub {
//	ctx := context.Background()
//	rdb := InitRedis()
//	subscriber := rdb.Subscribe(ctx, accessLogKey)
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
