package database

import (
	"github.com/pieterclaerhout/go-log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

type MariaDBConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

func (conf MariaDBConf) get() {
	conf.Host = os.Getenv("MARIADB_HOST")
	conf.User = os.Getenv("MARIADB_USER")
	conf.Pass = os.Getenv("MARIADB_PASS")
	conf.Data = os.Getenv("MARIADB_DATA")
}

func InitMariaDB() *gorm.DB {
	cfg := &MariaDBConf{}
	cfg.get()
	dsn := cfg.User + ":" + cfg.Pass + "@tcp(" + cfg.Host + ":3306)/" + cfg.Data
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return db
}
