package router

import (
	commonSvc "fatalisa-public-api/service/common"
	prayScheduleSvc "fatalisa-public-api/service/common/pray-schedule"
	qrisSvc "fatalisa-public-api/service/qris"
	"github.com/gin-gonic/gin"
)

func (router *Config) initLandingRoute() {
	router.Gin.GET("/", func(c *gin.Context) {
		c.SecureJSON(200, &commonSvc.Body{Message: "Welcome"})
	})
}

func (router *Config) initHealthRoute() {
	router.Gin.GET("/health", func(c *gin.Context) {
		c.SecureJSON(200, &commonSvc.Body{Message: "pong"})
	})
}

func (router *Config) versionChecker() {
	router.Gin.GET("/version", func(c *gin.Context) {
		c.SecureJSON(200, commonSvc.VersionChecker())
	})
}

func (router *Config) initApis() {
	api := router.Gin.Group("/api")
	{
		api.GET("/datetime", func(c *gin.Context) {
			c.SecureJSON(200, commonSvc.DateTimeApiService())
		})
		api.GET("/pray-schedule/city-list", func(c *gin.Context) {
			c.SecureJSON(200, prayScheduleSvc.GetCityList())
		})
		api.GET("/pray-schedule/:city", func(c *gin.Context) {
			c.SecureJSON(200, prayScheduleSvc.GetScheduleService(c))
		})
		api.POST("/pray-schedule", func(c *gin.Context) {
			c.SecureJSON(200, prayScheduleSvc.GetScheduleService(c))
		})
		qrisGroup := api.Group("/qris")
		{
			qrisGroup.GET("/mpm/:raw", func(c *gin.Context) {
				c.SecureJSON(200, qrisSvc.ParseMpmService(c))
			})
			qrisGroup.POST("/mpm", func(c *gin.Context) {
				c.SecureJSON(200, qrisSvc.ParseMpmService(c))
			})
			qrisGroup.POST("/cpm", func(c *gin.Context) {
				c.SecureJSON(200, qrisSvc.ParseCpmService(c))
			})
		}
	}
}
