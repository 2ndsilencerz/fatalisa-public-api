package common

import "github.com/pieterclaerhout/go-log"

func InitLogger() {
	log.DebugMode = true
	log.DebugSQLMode = true
	log.PrintTimestamp = true
	log.PrintColors = true
	log.TimeFormat = "2006-01-02 15:04:05 -0700"
}
