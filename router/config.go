package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
	"io"
	"os"
	"strconv"
)

type Config struct {
	Gin *gin.Engine
}

func ginLoggerTask(kind string, c *gin.Context) {
	if len(c.Request.RequestURI) > 0 && c.Request.RequestURI != "/health" {
		kindStr := fmt.Sprintf("%-10s", kind)
		reqMethod := fmt.Sprintf("%-10s", c.Request.Method)
		reqUri := fmt.Sprintf("%-10s ", c.Request.RequestURI)
		statusCode := fmt.Sprintf("%-3s ", strconv.Itoa(c.Writer.Status()))
		clientIP := fmt.Sprintf("%-16s ", c.ClientIP())
		hostHeader := c.Request.Header.Get("X-Real-Ip")
		if len(hostHeader) > 0 {
			clientIP = fmt.Sprintf("%-16s ", hostHeader)
		}
		log.Info(kindStr, clientIP, reqMethod, reqUri, statusCode)
	}
}

func ginCustomLogger(c *gin.Context) {
	ginLoggerTask("Request", c)
	c.Next()
	ginLoggerTask("Response", c)
}

func ginLogHandler() gin.HandlerFunc {
	return ginCustomLogger
}

func (router *Config) Get() {
	router.Gin = gin.New()
	router.Gin.Use(ginLogHandler())
	gin.DefaultWriter = io.MultiWriter(log.Default.Out)
	gin.ForceConsoleColor()
}

func (router *Config) Run() {
	router.Get()
	router.InitRoutes()
	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "80"
	}
	log.Info("Service running at port ", port)
	if err := router.Gin.Run(":" + port); err != nil {
		log.Error(err)
		panic(err)
	}
}

func (router *Config) InitRoutes() {
	router.initLandingRoute()
	router.initHealthRoute()
	router.versionChecker()
	router.initApis()
}
