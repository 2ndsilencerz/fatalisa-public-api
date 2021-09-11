package config

import "os"

type List struct {
	MongoDBHost string `json:"mongo_db_host"`
	MongoDBUser string `json:"mongo_db_user"`
	MongoDBPass string `json:"mongo_db_pass"`
	MongoDBData string `json:"mongo_db_data"`

	RedisHost string `json:"redis_host"`
	RedisPass string `json:"redis_pass"`

	ConsulHost string `json:"consul_host"`
	ConsulPort string `json:"consul_port"`
}

func (list List) Get() {
	list.MongoDBHost = os.Getenv("MONGODB_HOST")
	list.MongoDBUser = os.Getenv("MONGODB_USER")
	list.MongoDBPass = os.Getenv("MONGODB_PASS")
	list.MongoDBData = os.Getenv("MONGODB_DATA")

	list.RedisHost = os.Getenv("REDIS_HOST")
	list.RedisPass = os.Getenv("REDIS_PASS")

	list.ConsulHost = os.Getenv("CONSUL_HOST")
	list.ConsulPort = os.Getenv("CONSUL_PORT")
}
