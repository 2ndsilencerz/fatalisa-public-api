package database

import (
	"github.com/go-redis/redis/v8"
	"os"
)

type RedisConf struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}

var redisCfg *RedisConf
var HeaderRedis = "redis"

func (conf *RedisConf) Get() {
	conf.Host, _ = os.LookupEnv("REDIS_HOST")
	conf.Password, _ = os.LookupEnv("REDIS_PASS")
}

func init() {
	redisCfg = &RedisConf{}
	redisCfg.Get()
}

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host + ":6379",
		Password: redisCfg.Password,
		DB:       0,
	})
	return rdb
}
