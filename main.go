package main

import (
	"fatalisa-public-api/config"
	"fatalisa-public-api/router"
	"github.com/subchen/go-log"
)

func init() {
	config.Init()
}

func main() {
	// any code after routerInit.Run() won't be executed, place it in init_task.go inside config
	log.Info("Starting service")
	routerInit := &router.Config{}
	routerInit.Run()
}
