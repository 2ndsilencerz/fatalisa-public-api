package router

import (
	"fatalisa-public-api/service/common"
	praySchedule "fatalisa-public-api/service/pray-schedule"
	"fatalisa-public-api/service/qris"
	"fatalisa-public-api/service/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (router *Config) initLandingRoute() {
	router.Gin.GET("/", func(c *gin.Context) {
		web.Service(c, "index")
	})
}

func (router *Config) initHealthRoute() {
	router.Gin.GET("/health", func(c *gin.Context) {
		c.SecureJSON(http.StatusOK, &common.Body{Message: "pong"})
	})
}

func (router *Config) versionChecker() {
	router.Gin.GET("/version", func(c *gin.Context) {
		c.SecureJSON(http.StatusOK, common.VersionChecker())
	})
}

func (router *Config) initApis() {
	api := router.Gin.Group("/api")
	{
		api.GET("/datetime", func(c *gin.Context) {
			c.SecureJSON(http.StatusOK, common.DateTimeApi())
		})
		api.GET("/pray-schedule/city-list", func(c *gin.Context) {
			c.SecureJSON(http.StatusOK, praySchedule.GetCityList())
		})
		api.GET("/pray-schedule/:city", func(c *gin.Context) {
			c.SecureJSON(http.StatusOK, praySchedule.GetSchedule(c))
		})
		api.POST("/pray-schedule", func(c *gin.Context) {
			c.SecureJSON(http.StatusOK, praySchedule.GetSchedule(c))
		})
		qrisGroup := api.Group("/qris")
		{
			qrisGroup.GET("/mpm/:raw", func(c *gin.Context) {
				c.SecureJSON(http.StatusOK, qris.ParseMpm(c))
			})
			qrisGroup.POST("/mpm", func(c *gin.Context) {
				c.SecureJSON(http.StatusOK, qris.ParseMpm(c))
			})
			qrisGroup.POST("/cpm", func(c *gin.Context) {
				c.SecureJSON(http.StatusOK, qris.ParseCpm(c))
			})
		}
	}
}
