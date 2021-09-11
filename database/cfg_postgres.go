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

func (conf PostgresConf) get() {
	conf.Host = os.Getenv("POSTGRES_HOST")
	conf.User = os.Getenv("POSTGRES_USER")
	conf.Pass = os.Getenv("POSTGRES_PASS")
	conf.Data = os.Getenv("POSTGRES_DATA")
}

func InitPostgres() *gorm.DB {
	cfg := &PostgresConf{}
	cfg.get()
	dsn := "host=" + cfg.Host +
		" user=" + cfg.User +
		" password=" + cfg.Pass +
		" dbname=" + cfg.Data +
		" port=5432 TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return db
}
