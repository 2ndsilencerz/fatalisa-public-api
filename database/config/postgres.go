package config

import (
	"fatalisa-public-api/database"
	"github.com/subchen/go-log"
	"gorm.io/gorm"
	"os"
)

type PostgresConf struct {
	database.DBConf
	Client *gorm.DB `json:"client"`
}

var postgresCfg *PostgresConf

func (conf *PostgresConf) GetSettings() {
	conf.Host, _ = os.LookupEnv("POSTGRES_HOST")
	conf.User, _ = os.LookupEnv("POSTGRES_USER")
	conf.Pass, _ = os.LookupEnv("POSTGRES_PASS")
	conf.Data, _ = os.LookupEnv("POSTGRES_DATA")
}

func InitPostgres() *PostgresConf {
	if postgresCfg == nil || postgresCfg.Client == nil {
		postgresCfg = &PostgresConf{}
		postgresCfg.GetSettings()

		var db *gorm.DB
		//var err error
		//dsn := "host=" + postgresCfg.Host +
		//	" user=" + postgresCfg.User +
		//	" password=" + postgresCfg.Pass +
		//	" dbname=" + postgresCfg.Data +
		//	" port=5432 TimeZone=Asia/Jakarta"
		//if db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//	Logger: logger.Default.LogMode(logger.Silent),
		//}); err != nil {
		//	log.Error(err)
		//}
		postgresCfg.Client = db
	}
	return postgresCfg
}

func checkPostgres() {
	db := InitPostgres()
	if db.Client == nil {
		printConf(db)
		db = nil
	}
}

func (conf *PostgresConf) AutoMigrate(v interface{}) {
	if conf.Client != nil && v != nil {
		if err := conf.Client.AutoMigrate(v); err != nil {
			log.Error(err)
		}
	}
}

func (conf *PostgresConf) Write(v interface{}) {
	if conf.Client != nil && v != nil {
		if err := conf.Client.Create(&v); err != nil {
			log.Error(err)
		}
	}
}

//func (conf *PostgresConf) Select(v interface{}, query ...interface{}) {
//	if conf.Client != nil && v != nil {
//		conf.Client.Find(&v, query)
//	}
//}
