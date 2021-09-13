package database

import (
	"fmt"
	"github.com/pieterclaerhout/go-log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

type MariaDBConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

var mariaDbCfg *MariaDBConf

func (conf *MariaDBConf) Get() {
	conf.Host, _ = os.LookupEnv("MARIADB_HOST")
	conf.User, _ = os.LookupEnv("MARIADB_USER")
	conf.Pass, _ = os.LookupEnv("MARIADB_PASS")
	conf.Data, _ = os.LookupEnv("MARIADB_DATA")
}

func init() {
	mariaDbCfg = &MariaDBConf{}
	mariaDbCfg.Get()
}

func InitMariaDB() *gorm.DB {
	var db *gorm.DB
	dsn := mariaDbCfg.User + ":" + mariaDbCfg.Pass + "@tcp(" + mariaDbCfg.Host + ":3306)/" + mariaDbCfg.Data
	//log.Info(HeaderGorm, "|", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		str := fmt.Sprintf("%-8s", "mariadb")
		log.Error(str, "|", err)
		//panic(err)
	}
	return db
}
