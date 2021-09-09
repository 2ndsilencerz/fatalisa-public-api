package common

import (
	"fatalisa-public-api/utils"
	"github.com/gin-gonic/gin"
)

type Router struct {
	R *gin.Engine
}

func (router *Router) get() {
	router.R = gin.Default()
}

func (router *Router) Run() {
	router.get()
	router.initRoutes()
	err := router.R.Run(":80")
	ErrPrint(err, true)
}

func (router *Router) initRoutes() {
	router.initHealthRoute()
	router.initApis()
}

func (router *Router) initHealthRoute() {
	router.R.GET("/", func(c *gin.Context) {
		c.JSON(200, &utils.Body{Message: "Welcome"})
	})
	router.R.GET("/health", func(c *gin.Context) {
		c.JSON(200, &utils.Body{Message: "pong"})
	})
}

func (router *Router) initApis() {
	router.R.GET("/api/datetime", func(c *gin.Context) {
		c.JSON(200, utils.DatetimeApi())
	})
}
