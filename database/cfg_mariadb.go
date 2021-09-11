package database

import (
	"github.com/pieterclaerhout/go-log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log2 "log"
	"os"
)

type MariaDBConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

var mariaDbCfg *MariaDBConf

func (conf MariaDBConf) Get() {
	conf.Host = os.Getenv("MARIADB_HOST")
	conf.User = os.Getenv("MARIADB_USER")
	conf.Pass = os.Getenv("MARIADB_PASS")
	conf.Data = os.Getenv("MARIADB_DATA")
}

func init() {
	mariaDbCfg = &MariaDBConf{}
	mariaDbCfg.Get()
}

func InitMariaDB() *gorm.DB {
	dsn := mariaDbCfg.User + ":" + mariaDbCfg.Pass + "@tcp(" + mariaDbCfg.Host + ":3306)/" + mariaDbCfg.Data
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
		log2.Println(err)
		panic(err)
	}
	return db
}
