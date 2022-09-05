package common

import (
	"os"
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
	if buildDate, exist := os.LookupEnv("BUILD_DATE"); exist && len(buildDate) > 0 {
		res.Message = buildDate
		return &res
	}

	if buildDate, err := os.ReadFile("/build-date.txt"); err == nil {
		res.Message = string(buildDate)
		return &res
	}
	return &res
}
