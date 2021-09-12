package database

import (
	"github.com/gofrs/uuid"
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

func (accessLog AccessLog) WriteToMariaDB() {
	if db := InitMariaDB(); db != nil {
		db.Create(&accessLog)
		Close(db)
	}
}

func (accessLog AccessLog) WriteToPostgres() {
	if db := InitPostgres(); db != nil {
		db.Create(&accessLog)
		Close(db)
	}
}

func (accessLog AccessLog) WriteLog() {
	accessLog.WriteToMariaDB()
	accessLog.WriteToPostgres()
}
