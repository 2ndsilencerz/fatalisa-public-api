package database

/*
This function used to naming a table with custom name
Only used with Gorm.io framework
*/

type Tabler interface {
	TableName() string
}

func (AccessLog) TableName() string {
	return "access_log"
}

func (ErrorLog) TableName() string {
	return "error_log"
}
