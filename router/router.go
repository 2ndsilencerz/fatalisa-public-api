package router

import (
	"fatalisa-public-api/database/config"
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

var HeaderGin = fmt.Sprintf("%-8s", "gin")

func loggerTask(kind string, c *gin.Context) {
	kind = fmt.Sprintf("%-10s", kind)
	reqMethod := fmt.Sprintf("%-5s", c.Request.Method)
	reqUri := fmt.Sprintf("%s", c.Request.RequestURI)
	var statusCode int
	if kind == "Response" {
		statusCode = c.Writer.Status()
	}
	clientIP := fmt.Sprintf("%s", c.ClientIP())
	log.Info(HeaderGin, "|", kind, clientIP, reqMethod, reqUri, statusCode)
	go saveLogToDB(kind, c)
}

func saveLogToDB(kind string, ctxCopy *gin.Context) {
	accessLog := &entity.AccessLog{
		Kind:     kind,
		IP:       ctxCopy.ClientIP(),
		Method:   ctxCopy.Request.Method,
		FullPath: ctxCopy.FullPath(),
	}
	if kind == "Response" {
		accessLog.StatusCode = ctxCopy.Writer.Status()
	}
	//accessLog.WriteLog()
	config.PutToRedisQueue(&accessLog, entity.AccessLogKey)
}

func ginCustomLogger(c *gin.Context) {
	loggerTask("Request", c)
	c.Next()
	loggerTask("Response", c)
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
	log.Info(HeaderGin, "|", "Service running", port)
	if err := router.Gin.Run(":" + port); err != nil {
		log.Error(HeaderGin, "|", err)
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
			response := commonSvc.DatetimeApi()
			c.SecureJSON(200, &response)
		})
		api.GET("/pray-schedule/city-list", func(c *gin.Context) {
			c.SecureJSON(200, pray_schedule.GetCityList())
		})
		api.POST("/pray-schedule", func(c *gin.Context) {
			jsonBody := &pray_schedule.PrayScheduleReq{}
			if err := c.BindJSON(jsonBody); err != nil {
				log.Error(pray_schedule.HeaderPray, "|", err)
			} else {
				response := pray_schedule.GetSchedule(jsonBody)
				c.SecureJSON(200, &response)
			}
		})
		qrisGroup := api.Group("/qris")
		{
			qrisGroup.GET("/mpm/:raw", func(c *gin.Context) {
				raw := c.Param("raw")
				result := qrisSvc.MpmData{}
				result.GetData(raw)
				c.SecureJSON(200, &result)
			})
			qrisGroup.POST("/mpm", func(c *gin.Context) {
				mpmReq := &qrisSvc.MpmRequest{}
				if err := c.BindJSON(mpmReq); err != nil {
					log.Error(qrisSvc.HeaderMpm, "|", err)
				} else {
					result := qrisSvc.MpmData{}
					result.GetData(mpmReq.Raw)
					c.SecureJSON(200, &result)
				}
			})
			qrisGroup.POST("/cpm", func(c *gin.Context) {
				cpmReq := &qrisSvc.CpmRequest{}
				if err := c.BindJSON(cpmReq); err != nil {
					log.Error(qrisSvc.HeaderCpm, "|", err)
				} else {
					result := qrisSvc.CpmData{}
					result.GetData(cpmReq.Raw)
					c.SecureJSON(200, &result)
				}
			})
		}
	}
}
