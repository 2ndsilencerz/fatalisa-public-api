package utils

import (
	"encoding/json"
	"github.com/pieterclaerhout/go-log"
)

func Jsonify(v interface{}) string {
	var j []byte
	var err error
	if j, err = json.Marshal(v); err != nil {
		log.Error(err)
	}
	return string(j)
}
