package config

import (
	"context"
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
	var err error
	ctx := context.Background()
	dsn := "mongodb://" + mongoDbConf.User + ":" + mongoDbConf.Pass + "@" + mongoDbConf.Host + ":27017/"
	if db, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn)); err != nil {
		log.Error(err)
	}
	return db, ctx, mongoDbConf
}
