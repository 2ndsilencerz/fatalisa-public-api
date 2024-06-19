package config

import (
	"context"
	"encoding/json"
	"fatalisa-public-api/database"
	utils2 "fatalisa-public-api/service/web/utils"
	"fatalisa-public-api/utils"
	"github.com/go-redis/redis/v8"
	"os"
	"time"
)

type RedisConf struct {
	database.DBConf
	Client *redis.Client `json:"client"`
}

var redisCfg *RedisConf

func (conf *RedisConf) GetSettings() {
	conf.Host, _ = os.LookupEnv("REDIS_HOST")
	conf.Pass, _ = os.LookupEnv("REDIS_PASS")
}

func InitRedis() *RedisConf {
	if redisCfg == nil || redisCfg.Client == nil {
		redisCfg = &RedisConf{}
		redisCfg.GetSettings()

		if len(redisCfg.Host) > 0 {
			rdb := redis.NewClient(&redis.Options{
				Addr:     redisCfg.Host + ":6379",
				Password: redisCfg.Pass,
				DB:       0,
			})
			redisCfg.Client = rdb
		}
	}
	return redisCfg
}

//func checkRedis() {
//	ctx := context.Background()
//	rdb := InitRedis()
//	if !rdb.connected(ctx) {
//		printConf(rdb)
//		rdb = nil
//	}
//}

func (conf *RedisConf) PushQueue(key string, v interface{}, ctx context.Context) {
	//ctx := context.Background()
	rdb := InitRedis()
	if conf.connected(ctx) && v != nil {
		rawString := utils.Jsonify(v)
		if len(rawString) > 0 {
			utils2.ErrorHandler(rdb.Client.LPush(ctx, key, rawString).Err())
		}
	}
}

func (conf *RedisConf) PopQueue(key string, v interface{}, ctx context.Context) {
	rdb := InitRedis()
	//ctx := context.Background()
	if conf.connected(ctx) && v != nil {
		rawString := rdb.Client.RPop(ctx, key).Val()
		if len(rawString) > 0 {
			utils2.ErrorHandler(json.Unmarshal([]byte(rawString), v))
		}
	}
	v = nil
}

//func (conf *RedisConf) SetString(key string, value string, ctx context.Context) bool {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	added := false
//	if rdb.connected(ctx) && len(value) > 0 {
//		err := rdb.Client.Set(ctx, key, value, 0).Err()
//		if err, _ := utils2.ErrorHandler(err); !err {
//			added = true
//		}
//	}
//	return added
//}

func (conf *RedisConf) Set(key string, value interface{}, ctx context.Context, duration time.Duration) bool {
	rdb := InitRedis()
	//ctx := context.Background()
	added := false
	if rdb.connected(ctx) && value != nil {
		err := rdb.Client.Set(ctx, key, value, duration).Err()
		if err, _ := utils2.ErrorHandler(err); !err {
			added = true
		}
	}
	return added
}

//func (conf *RedisConf) SetHash(key string, value interface{}, ctx context.Context, duration time.Duration) bool {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	added := false
//	if rdb.connected(ctx) && value != nil {
//		jsonForm, err := json.Marshal(value)
//		if err, _ := utils2.ErrorHandler(err); err {
//			return false
//		}
//		err = rdb.Client.HSet(ctx, key, jsonForm, duration).Err()
//		if err, _ := utils2.ErrorHandler(err); !err {
//			added = true
//		}
//	}
//	return added
//}
//
//func (conf *RedisConf) SetMultiHash(key string, value interface{}, ctx context.Context, duration time.Duration) bool {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	added := false
//	if rdb.connected(ctx) && value != nil {
//		jsonForm, err := json.Marshal(value)
//		if err, _ := utils2.ErrorHandler(err); err {
//			return false
//		}
//		err = rdb.Client.HMSet(ctx, key, jsonForm, duration).Err()
//		if err, _ := utils2.ErrorHandler(err); !err {
//			added = true
//		}
//	}
//	return added
//}
//
//func (conf *RedisConf) GetString(key string, ctx context.Context) string {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	result := ""
//	if rdb.connected(ctx) {
//		cmd := rdb.Client.Get(ctx, key)
//		if err, _ := utils2.ErrorHandler(cmd.Err()); !err {
//			result = cmd.Val()
//		}
//	}
//	return result
//}

func (conf *RedisConf) Get(key string, ctx context.Context) interface{} {
	rdb := InitRedis()
	//ctx := context.Background()
	var result interface{}
	if rdb.connected(ctx) {
		cmd := rdb.Client.Get(ctx, key)
		if err, _ := utils2.ErrorHandler(cmd.Err()); !err {
			result = cmd.Val()
		}
	}
	return result
}

//func (conf *RedisConf) GetHash(key, field string, ctx context.Context) interface{} {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	var result interface{}
//	if rdb.connected(ctx) {
//		cmd := rdb.Client.HGet(ctx, key, field)
//		if err, _ := utils2.ErrorHandler(cmd.Err()); !err {
//			result = cmd.Val()
//		}
//	}
//	return result
//}
//
//func (conf *RedisConf) GetMultiHash(key string, field []string, ctx context.Context) []string {
//	rdb := InitRedis()
//	//ctx := context.Background()
//	var result []string
//	if rdb.connected(ctx) {
//		//log.Info(key, field)
//		cmd := rdb.Client.HMGet(ctx, key, field...)
//		err := cmd.Err()
//		utils2.ErrorHandler(err)
//		log.Info(cmd.Val())
//		err = json.Unmarshal([]byte(fmt.Sprint(cmd.Val())), &result)
//		utils2.ErrorHandler(err)
//	}
//	log.Info(result)
//	return result
//}

func (conf *RedisConf) GetKeys(pattern string, ctx context.Context) map[string]string {
	rdb := InitRedis()
	//ctx := context.Background()
	var result = make(map[string]string)
	if rdb.connected(ctx) {
		cmd := rdb.Client.Keys(ctx, pattern)
		if err, _ := utils2.ErrorHandler(cmd.Err()); !err {
			listResult := cmd.Val()
			if len(listResult) > 0 {
				for _, value := range listResult {
					newCmd := rdb.Client.Get(ctx, value)
					if newV, _ := utils2.ErrorHandler(newCmd.Err()); !newV {
						result[value] = newCmd.Val()
					}
				}
			}
		}
	}
	return result
}

func (conf *RedisConf) connected(ctx context.Context) bool {
	if len(conf.Host) > 0 {
		if err, _ := utils2.ErrorHandler(conf.Client.Ping(ctx).Err()); !err {
			return true
		}
	}
	return false
}
