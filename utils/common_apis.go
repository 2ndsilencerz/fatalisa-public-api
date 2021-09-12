package utils

import "time"

type Body struct {
	Message string `json:"message"`
}

func DatetimeApi() *Body {
	currentTime := time.Now().Format("2006-01-02 15:04:05 -0700")
	return &Body{currentTime}
}
