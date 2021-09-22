package common

type Tabler interface {
	TableName() string
}

func (ErrorLog) TableName() string {
	return errorLogKey
}
