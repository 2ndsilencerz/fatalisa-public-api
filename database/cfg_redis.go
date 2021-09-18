package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pieterclaerhout/go-log"
	"os"
)

type RedisConf struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}

var redisCfg *RedisConf
var HeaderRedis = fmt.Sprintf("%-8s", "redis")

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

func PutToRedisQueue(v interface{}, key string) {
	if rawString, err := json.Marshal(v); err != nil {
		log.Error(key, "|", err)
	} else if rdb := InitRedis(); rdb != nil {
		defer CloseRedis(rdb)
		ctx := context.Background()
		if errorPush := rdb.LPush(ctx, key, string(rawString)).Err(); errorPush != nil {
			log.Error(HeaderRedis, "|", errorPush)
		}
	}
}
