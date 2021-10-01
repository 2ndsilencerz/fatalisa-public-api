package config

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/pieterclaerhout/go-log"
	"os"
)

type RedisConf struct {
	Host     string        `json:"host"`
	Password string        `json:"password"`
	Client   *redis.Client `json:"client"`
}

var redisCfg *RedisConf

func (conf *RedisConf) Get() {
	conf.Host, _ = os.LookupEnv("REDIS_HOST")
	conf.Password, _ = os.LookupEnv("REDIS_PASS")
}

func InitRedis() *RedisConf {
	if redisCfg == nil || redisCfg.Client == nil {
		redisCfg = &RedisConf{}
		redisCfg.Get()

		rdb := redis.NewClient(&redis.Options{
			Addr:     redisCfg.Host + ":6379",
			Password: redisCfg.Password,
			DB:       0,
		})
		redisCfg.Client = rdb
	}
	return redisCfg
}

func checkRedis() {
	ctx := context.Background()
	rdb := InitRedis()
	if err := rdb.Client.Ping(ctx).Err(); err != nil {
		printConf(rdb)
		rdb = nil
	}
}

func (conf *RedisConf) PushQueue(key string, v interface{}) {
	ctx := context.Background()
	rdb := InitRedis()
	if err := rdb.Client.Ping(ctx).Err(); err != nil {
		log.Error(err)
	} else if v != nil {
		rawString := utils.Jsonify(v)
		if len(rawString) > 0 {
			if errorPush := rdb.Client.LPush(ctx, key, rawString).Err(); errorPush != nil {
				log.Error(errorPush)
			}
		}
	}
}

func (conf *RedisConf) PopQueue(key string, v interface{}) {
	rdb := InitRedis()
	ctx := context.Background()
	if err := rdb.Client.Ping(ctx).Err(); err != nil {
		log.Error(err)
	} else if v != nil {
		ctx := context.Background()
		rawString := rdb.Client.RPop(ctx, key).Val()
		if len(rawString) > 0 {
			if err := json.Unmarshal([]byte(rawString), v); err != nil {
				log.Error(err)
			}
		}
	}
	v = nil
}
