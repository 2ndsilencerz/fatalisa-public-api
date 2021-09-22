package main

import (
	"fatalisa-public-api/config"
	"fatalisa-public-api/router"
	"github.com/pieterclaerhout/go-log"
)

func init() {
	config.Init()
}

func main() {
	// any code after routerInit.Run() won't be executed, place it in init() above
	log.Info("Starting service")
	routerInit := &router.Config{}
	routerInit.Run()
}
