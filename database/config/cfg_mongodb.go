package config

import (
	"context"
	"fmt"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type MongoDBConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

var mongoDbConf *MongoDBConf
var HeaderMongoDB = fmt.Sprintf("%-8s", "mongodb")

func (conf *MongoDBConf) Get() {
	conf.Host, _ = os.LookupEnv("MONGODB_HOST")
	conf.User, _ = os.LookupEnv("MONGODB_USER")
	conf.Pass, _ = os.LookupEnv("MONGODB_PASS")
	conf.Data, _ = os.LookupEnv("MONGODB_DATA")
}

func init() {
	mongoDbConf = &MongoDBConf{}
	mongoDbConf.Get()
}

func InitMongoDB() (*mongo.Client, context.Context, *MongoDBConf) {
	var db *mongo.Client
	ctx := context.Background()
	dsn := "mongodb://" + mongoDbConf.User + ":" + mongoDbConf.Pass + "@" + mongoDbConf.Host + ":27017/"
	//log.Info(HeaderGorm, "|", dsn)
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		log.Error(HeaderMongoDB, "|", err)
	}
	return db, ctx, mongoDbConf
}
