package common

import (
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"io/ioutil"
	"time"
)

type Body struct {
	Message string `json:"message"`
}

func DateTimeApiService() *Body {
	body := datetimeApi()
	log.Info(utils.Jsonify(body))
	return body
}

func datetimeApi() *Body {
	currentTime := time.Now().Format("2006/01/02 15:04:05 -0700")
	return &Body{Message: currentTime}
}

func VersionChecker() *Body {
	res := Body{}
	if version, err := ioutil.ReadFile("/build-date.txt"); err != nil {
		log.Error(err)
	} else {
		res.Message = string(version)
	}
	return &res
}
