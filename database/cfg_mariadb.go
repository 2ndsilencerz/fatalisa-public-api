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
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		str := fmt.Sprintf("%-8s", "mariadb")
		log.Error(str, "|", err)
		//panic(err)
	} else {
		return db
	}
	return nil
}
