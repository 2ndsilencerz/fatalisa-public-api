package utils

import (
	"encoding/json"
	"github.com/pieterclaerhout/go-log"
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

func GetDuration(duration string) time.Duration {
	var res time.Duration
	var err error
	res, err = time.ParseDuration(duration)
	if err != nil {
		log.Error(err)
	}
	return res
}
