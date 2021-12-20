package database

type DBConf struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Data string `json:"data"`
}

type DBInterface interface {
	GetSettings()
}

type GormDBInterface interface {
	DBInterface
	Write(v interface{})
	AutoMigrate(v interface{})
}

type MongoDBInterface interface {
	DBInterface
	InsertOne(collection string, v interface{})
}

type RedisInterface interface {
	DBInterface
	Get(key string) string
	Set(key string, value string) bool
	PushQueue(key string, v interface{})
	PopQueue(key string, v interface{})
}
