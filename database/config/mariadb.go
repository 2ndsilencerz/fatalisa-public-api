package config

import (
	"fatalisa-public-api/database"
	"github.com/subchen/go-log"
	"gorm.io/gorm"
	"os"
)

type MariaDBConf struct {
	database.DBConf
	Client *gorm.DB `json:"client"`
}

var mariaDbCfg *MariaDBConf

func (conf *MariaDBConf) GetSettings() {
	conf.Host, _ = os.LookupEnv("MARIADB_HOST")
	conf.User, _ = os.LookupEnv("MARIADB_USER")
	conf.Pass, _ = os.LookupEnv("MARIADB_PASS")
	conf.Data, _ = os.LookupEnv("MARIADB_DATA")
}

func InitMariaDB() *MariaDBConf {
	if mariaDbCfg == nil || mariaDbCfg.Client == nil {
		mariaDbCfg = &MariaDBConf{}
		mariaDbCfg.GetSettings()

		var db *gorm.DB
		//var err error
		//dsn := mariaDbCfg.User + ":" + mariaDbCfg.Pass + "@tcp(" + mariaDbCfg.Host + ":3306)/" + mariaDbCfg.Data
		//if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//	Logger: logger.Default.LogMode(logger.Silent),
		//}); err != nil {
		//	log.Error(err)
		//}
		mariaDbCfg.Client = db
	}
	return mariaDbCfg
}

func checkMariaDB() {
	db := InitMariaDB()
	if db.Client == nil {
		printConf(db)
		db = nil
	}
}

func (conf *MariaDBConf) AutoMigrate(v interface{}) {
	if conf.Client != nil && v != nil {
		if err := conf.Client.AutoMigrate(v); err != nil {
			log.Error(err)
		}
	}
}

func (conf *MariaDBConf) Write(v interface{}) {
	if conf.Client != nil && v != nil {
		if err := conf.Client.Create(&v); err != nil {
			log.Error(err)
		}
	}
}

//func (conf *MariaDBConf) Select(v interface{}, query string) {
//	if conf.Client != nil && v != nil {
//		conf.Client.Select(&v, query)
//	}
//}
