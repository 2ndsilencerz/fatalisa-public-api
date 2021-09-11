package database

import (
	"github.com/pieterclaerhout/go-log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type PostgresConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

var postgresCfg *PostgresConf

func (conf PostgresConf) Get() {
	conf.Host = os.Getenv("POSTGRES_HOST")
	conf.User = os.Getenv("POSTGRES_USER")
	conf.Pass = os.Getenv("POSTGRES_PASS")
	conf.Data = os.Getenv("POSTGRES_DATA")
}

func init() {
	postgresCfg = &PostgresConf{}
	postgresCfg.Get()
}

func InitPostgres() *gorm.DB {
	dsn := "host=" + postgresCfg.Host +
		" user=" + postgresCfg.User +
		" password=" + postgresCfg.Pass +
		" dbname=" + postgresCfg.Data +
		" port=5432 TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
		//panic(err)
	}
	return db
}
