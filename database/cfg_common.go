package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pieterclaerhout/go-log"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"time"
)

var HeaderGorm = fmt.Sprintf("%-8s", "gorm")

func Close(database *gorm.DB) {
	if db, err := database.DB(); err != nil {
		log.Error(HeaderGorm, "|", err)
		panic(err)
	} else {
		if err = db.Close(); err != nil {
			log.Error(HeaderGorm, "|", err)
			panic(err)
		}
	}
}

func CloseMongo(database *mongo.Client, ctx context.Context) {
	if err := database.Disconnect(ctx); err != nil {
		log.Error(HeaderMongoDB, "|", err)
		panic(err)
	}
}

func CloseRedis(client *redis.Client) {
	if err := client.Close(); err != nil {
		log.Error(HeaderRedis, "|", err)
		panic(err)
	}
}

func DbConnCheck() {
	for {
		log.Info(HeaderGorm, "|", "Checking DB connection")
		go checkPostgres()
		go checkMariaDB()
		go checkMongoDB()
		if sleepTime, err := time.ParseDuration("30s"); err != nil {
			log.Error(HeaderGorm, "|", err)
		} else {
			time.Sleep(sleepTime)
		}
	}
}

func checkPostgres() {
	if postgres := InitPostgres(); postgres != nil {
		Close(postgres)
	}
}

func checkMariaDB() {
	if mariadb := InitMariaDB(); mariadb != nil {
		Close(mariadb)
	}
}

func checkMongoDB() {
	if mongodb, ctx, _ := InitMongoDB(); mongodb != nil {
		CloseMongo(mongodb, ctx)
	}
}
