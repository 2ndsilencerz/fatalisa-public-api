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

func (list *List) Get() {
	list.MongoDBHost, _ = os.LookupEnv("MONGODB_HOST")
	list.MongoDBUser, _ = os.LookupEnv("MONGODB_USER")
	list.MongoDBPass, _ = os.LookupEnv("MONGODB_PASS")
	list.MongoDBData, _ = os.LookupEnv("MONGODB_DATA")

	list.RedisHost, _ = os.LookupEnv("REDIS_HOST")
	list.RedisPass, _ = os.LookupEnv("REDIS_PASS")

	list.ConsulHost, _ = os.LookupEnv("CONSUL_HOST")
	list.ConsulPort, _ = os.LookupEnv("CONSUL_PORT")
}
