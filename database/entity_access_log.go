package database

import (
	"github.com/gofrs/uuid"
	"github.com/pieterclaerhout/go-log"
	"gorm.io/gorm"
)

/*
Column name convention used in Gorm.io framework is snake_case
In-case you want to assign with different name, use tag `gorm:"column:columnName"`
*/

type AccessLog struct {
	gorm.Model
	UUID       uuid.UUID `json:"uuid" gorm:"column:uuid"`
	Kind       string    `json:"message" gorm:"column:kind"`
	IP         string    `json:"ip"`
	Method     string    `json:"method"`
	FullPath   string    `json:"full_path"`
	StatusCode int       `json:"status_code"`
	Created    int64     `gorm:"autoCreateTime,column:created" json:"created"`
}

func (accessLog *AccessLog) WriteToMariaDB() {
	if db := InitMariaDB(); db != nil {
		if err := db.AutoMigrate(&accessLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&accessLog)
		Close(db)
	}
}

func (accessLog *AccessLog) WriteToPostgres() {
	if db := InitPostgres(); db != nil {
		if err := db.AutoMigrate(&accessLog); err != nil {
			log.Error(HeaderGorm, "|", err)
		}
		db.Create(&accessLog)
		Close(db)
	}
}

func (accessLog *AccessLog) WriteLog() {
	uuidGenerated, err := uuid.NewV4()
	if err != nil {
		log.Error(HeaderGorm, "|", err)
	}
	accessLog.UUID = uuidGenerated
	accessLog.WriteToMariaDB()
	accessLog.WriteToPostgres()
}
