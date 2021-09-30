package config

import (
	"github.com/pieterclaerhout/go-log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type PostgresConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

var postgresCfg *PostgresConf

func (conf *PostgresConf) Get() {
	conf.Host, _ = os.LookupEnv("POSTGRES_HOST")
	conf.User, _ = os.LookupEnv("POSTGRES_USER")
	conf.Pass, _ = os.LookupEnv("POSTGRES_PASS")
	conf.Data, _ = os.LookupEnv("POSTGRES_DATA")
}

func init() {
	postgresCfg = &PostgresConf{}
	postgresCfg.Get()
}

func InitPostgres() *gorm.DB {
	var db *gorm.DB
	var err error
	dsn := "host=" + postgresCfg.Host +
		" user=" + postgresCfg.User +
		" password=" + postgresCfg.Pass +
		" dbname=" + postgresCfg.Data +
		" port=5432 TimeZone=Asia/Jakarta"
	if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		log.Error(err)
	}
	return db
}