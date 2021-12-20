package config

import (
	"context"
	"fatalisa-public-api/database"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type MongoDBConf struct {
	database.DBConf
	Client  *mongo.Client   `json:"client"`
	Context context.Context `json:"context"`
}

var mongoDbConf *MongoDBConf

func (conf *MongoDBConf) GetSettings() {
	conf.Host, _ = os.LookupEnv("MONGODB_HOST")
	conf.User, _ = os.LookupEnv("MONGODB_USER")
	conf.Pass, _ = os.LookupEnv("MONGODB_PASS")
	conf.Data, _ = os.LookupEnv("MONGODB_DATA")
	log.Info(utils.Jsonify(&conf))
}

func InitMongoDB() *MongoDBConf {
	if mongoDbConf == nil || mongoDbConf.Client == nil {
		mongoDbConf = &MongoDBConf{}
		mongoDbConf.GetSettings()

		var db *mongo.Client
		var err error
		ctx := context.Background()
		dsn := "mongodb://" +
			mongoDbConf.User + ":" +
			mongoDbConf.Pass + "@" +
			mongoDbConf.Host + ":27017/" +
			mongoDbConf.Data
		//+ "?authSource=" +
		//	mongoDbConf.Data + "directConnection=true&ssl=false"
		if db, err = mongo.Connect(ctx, options.Client().ApplyURI(dsn)); err != nil {
			log.Error(err)
		}
		mongoDbConf.Client = db
		mongoDbConf.Context = ctx
	}
	return mongoDbConf
}

func checkMongoDB() {
	db := InitMongoDB()
	if db.Client == nil {
		printConf(db)
		db = nil
	}
}

func (conf *MongoDBConf) InsertOne(collection string, v interface{}) {
	if conf.Client != nil && v != nil {

		selectedCollection := conf.Client.Database(conf.Data).Collection(collection)
		if bsonData, err := bson.Marshal(v); err != nil {
			log.Error(err)
		} else if _, err := selectedCollection.InsertOne(conf.Context, bsonData); err != nil {
			log.Error(err)
		}
	}
}
