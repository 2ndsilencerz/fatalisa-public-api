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
	UUID         uuid.UUID `json:"uuid" gorm:"column:uuid"`
	Message      string    `json:"message" gorm:"column:message"`
	IP           string    `json:"ip"`
	Method       string    `json:"method"`
	FullPath     string    `json:"full_path"`
	ResponseCode int       `json:"response_code"`
	PostData     string    `json:"post_data"`
	ParamData    string    `json:"param_data"`
	Response     string    `json:"response"`
	Created      int64     `gorm:"autoCreateTime,column:created" json:"created"`
}

func (accessLog AccessLog) WriteToMariaDB() {
	db := InitMariaDB()
	db.Create(&accessLog)
	Close(db)
}

func (accessLog AccessLog) WriteToPostgres() {
	db := InitPostgres()
	db.Create(&accessLog)
	Close(db)
}

func (accessLog AccessLog) WriteLog() {
	accessLog.WriteToMariaDB()
	accessLog.WriteToPostgres()
}
