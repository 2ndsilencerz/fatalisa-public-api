package utils

import (
	"encoding/json"
	"github.com/subchen/go-log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func Jsonify(v interface{}) string {
	var j []byte
	var err error
	if j, err = json.Marshal(v); err != nil {
		log.Error(err)
	}
	return string(j)
}

//func GetDuration(duration string) time.Duration {
//	var res time.Duration
//	var err error
//	res, err = time.ParseDuration(duration)
//	if err != nil {
//		log.Error(err)
//	}
//	return res
//}

func GetPodName() string {
	str := ""
	exist := false
	if str, exist = os.LookupEnv("POD_NAME"); !exist {
		rand.Seed(time.Now().UnixNano())
		str = time.Now().Format("2006-01-02") + "-" + strconv.Itoa(rand.Int())
	}
	return str
}

func Mkdir(location string) {
	err := os.Mkdir(location, os.FileMode(777))
	if err != nil {
		log.Error(err)
	}
}
