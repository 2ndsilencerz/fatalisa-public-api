package common

import (
	"github.com/subchen/go-log"
	"io/ioutil"
	"time"
)

type Body struct {
	Message string `json:"message"`
}

func DateTimeApi() *Body {
	currentTime := time.Now().Format("2006/01/02 15:04:05 -0700")
	return &Body{Message: currentTime}
}

func VersionChecker() *Body {
	res := Body{}
	if version, err := ioutil.ReadFile("/build-date.txt"); err != nil {
		log.Warn("No build date set!")
	} else {
		res.Message = string(version)
	}
	return &res
}
