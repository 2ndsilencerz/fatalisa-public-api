package router

import (
	"fatalisa-public-api/database/entity"
	commonSvc "fatalisa-public-api/service/common"
	"fatalisa-public-api/service/common/pray-schedule"
	qrisSvc "fatalisa-public-api/service/qris"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"os"
)

type Config struct {
	Gin *gin.Engine
}

func loggerTask(kind string, c *gin.Context) {
	kindStr := fmt.Sprintf("%-10s", kind)
	reqMethod := fmt.Sprintf("%-5s", c.Request.Method)
	reqUri := fmt.Sprintf("%s", c.Request.RequestURI)
	statusCode := c.Writer.Status()
	clientIP := fmt.Sprintf("%s", c.ClientIP())
	log.Info(kindStr, clientIP, reqMethod, reqUri, statusCode)
	go saveLogToDB(kind, c.Copy())
}

func saveLogToDB(kind string, c *gin.Context) {
	accessLog := &entity.AccessLog{
		Kind:       kind,
		IP:         c.ClientIP(),
		Method:     c.Request.Method,
		FullPath:   c.FullPath(),
		StatusCode: c.Writer.Status(),
	}
	accessLog.PutToRedisQueue()
}

func ginCustomLogger(c *gin.Context) {
	loggerTask("Request", c.Copy())
	c.Next()
	loggerTask("Response", c.Copy())
}

func ginLogHandler() gin.HandlerFunc {
	return ginCustomLogger
}

func (router *Config) Get() {
	router.Gin = gin.New()
	router.Gin.Use(ginLogHandler())
	gin.ForceConsoleColor()
}

func (router *Config) Run() {
	router.Get()
	router.InitRoutes()
	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "80"
	}
	log.Info("Service running", port)
	if err := router.Gin.Run(":" + port); err != nil {
		log.Error(err)
		panic(err)
	}
}

func (router *Config) InitRoutes() {
	router.initLandingRoute()
	router.initHealthRoute()
	router.initApis()
}

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

func (router *Config) initApis() {
	api := router.Gin.Group("/api")
	{
		api.GET("/datetime", func(c *gin.Context) {
			c.SecureJSON(200, commonSvc.DateTimeApiService())
		})
		api.GET("/pray-schedule/city-list", func(c *gin.Context) {
			c.SecureJSON(200, pray_schedule.GetCityList())
		})
		api.POST("/pray-schedule", func(c *gin.Context) {
			c.SecureJSON(200, pray_schedule.GetScheduleService(c))
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
