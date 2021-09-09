package config

import "fatalisa-public-api/common"

func Init() {
	initLogger()
	initRouter()
}

func initLogger() {
	common.InitLogger()
}

func initRouter() {
	router := &common.Router{}
	router.Run()
}
