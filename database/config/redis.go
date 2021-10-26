package config

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/utils"
	"github.com/go-redis/redis/v8"
	"github.com/subchen/go-log"
	"os"
)

type RedisConf struct {
	Host     string        `json:"host"`
	Password string        `json:"password"`
	Client   *redis.Client `json:"client"`
}

var redisCfg *RedisConf

func (conf *RedisConf) GetConfig() {
	conf.Host, _ = os.LookupEnv("REDIS_HOST")
	conf.Password, _ = os.LookupEnv("REDIS_PASS")
}

func InitRedis() *RedisConf {
	if redisCfg == nil || redisCfg.Client == nil {
		redisCfg = &RedisConf{}
		redisCfg.GetConfig()

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
	if !rdb.connected(ctx) {
		printConf(rdb)
		rdb = nil
	}
}

func (conf *RedisConf) PushQueue(key string, v interface{}) {
	ctx := context.Background()
	rdb := InitRedis()
	if conf.connected(ctx) && v != nil {
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
	if conf.connected(ctx) && v != nil {
		rawString := rdb.Client.RPop(ctx, key).Val()
		if len(rawString) > 0 {
			if err := json.Unmarshal([]byte(rawString), v); err != nil {
				log.Error(err)
			}
		}
	}
	v = nil
}

func (conf *RedisConf) Set(key string, value string) bool {
	rdb := InitRedis()
	ctx := context.Background()
	added := false
	if rdb.connected(ctx) && len(value) > 0 {
		err := rdb.Client.Set(ctx, key, value, 0).Err()
		if err != nil {
			log.Error(err)
		} else {
			added = true
		}
	}
	return added
}

func (conf *RedisConf) Get(key string) string {
	rdb := InitRedis()
	ctx := context.Background()
	result := ""
	if rdb.connected(ctx) {
		cmd := rdb.Client.Get(ctx, key)
		err := cmd.Err()
		if err != nil {
			log.Error(err)
		} else {
			result = cmd.String()
		}
	}
	return result
}

func (conf *RedisConf) connected(ctx context.Context) bool {
	if err := conf.Client.Ping(ctx).Err(); err != nil {
		log.Error(err)
		return false
	}
	return true
}
