package config

import (
	"context"
	"fatalisa-public-api/service/utils"
	"github.com/go-redis/redis/v8"
	"github.com/pieterclaerhout/go-log"
	"os"
)

type RedisConf struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}

var redisCfg *RedisConf

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
	rawString := utils.Jsonify(v)
	rdb := InitRedis()
	if len(rawString) > 0 && rdb != nil {
		defer CloseRedis(rdb)
		ctx := context.Background()
		if errorPush := rdb.LPush(ctx, key, rawString).Err(); errorPush != nil {
			log.Error(errorPush)
		}
	}
}
