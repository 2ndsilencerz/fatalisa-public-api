package qris

type Tabler interface {
	TableName() string
}

func (Log) TableName() string {
	return qrisKey
}
