package common

import (
	"fatalisa-public-api/utils"
	"github.com/pieterclaerhout/go-log"
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
