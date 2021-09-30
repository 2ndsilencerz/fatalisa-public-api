package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"os"
	"time"
)

type Config struct {
	Gin *gin.Engine
}

func loggerTask(kind string, c *gin.Context) {
	if len(c.Request.RequestURI) > 0 && c.Request.RequestURI != "/health" {
		kindStr := fmt.Sprintf("%-10s", kind)
		reqMethod := fmt.Sprintf("%-5s", c.Request.Method)
		reqUri := fmt.Sprintf("%s", c.Request.RequestURI)
		statusCode := c.Writer.Status()
		clientIP := fmt.Sprintf("%s", c.ClientIP())
		// hostHeader useful if we are using nginx port forwarding with following config
		// proxy_set_header Host $host;
		hostHeader := c.Request.Header.Get("Host")
		if len(hostHeader) > 0 {
			clientIP = hostHeader
		}
		log.Info(kindStr, clientIP, reqMethod, reqUri, statusCode)
		go saveLogToDB(kind, c.Copy())
	}
}

func saveLogToDB(kind string, c *gin.Context) {
	accessLog := &AccessLog{
		Kind:       kind,
		IP:         c.ClientIP(),
		Method:     c.Request.Method,
		FullPath:   c.Request.RequestURI,
		StatusCode: c.Writer.Status(),
		Created:    time.Now(),
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
