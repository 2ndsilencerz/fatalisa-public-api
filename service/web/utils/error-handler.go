package utils

import (
	"github.com/subchen/go-log"
)

func ErrorHandler(err error) (bool, string) {
	if err != nil {
		log.Errorf("Error: %v", err.Error())
		return true, err.Error()
	}
	return false, ""
}
