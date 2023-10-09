package utils

import "github.com/subchen/go-log"

func ErrorHandler(err error) (bool, string) {
	if err != nil {
		log.Errorf("Error: %v", err)
		return false, err.Error()
	}
	return true, ""
}
