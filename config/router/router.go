package router

import (
	"fatalisa-public-api/database"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"os"
)

type Router struct {
	R *gin.Engine
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
	saveLogToDB(kind, c)
}

func saveLogToDB(kind string, ctxCopy *gin.Context) {
	accessLog := &database.AccessLog{
		Kind:     kind,
		IP:       ctxCopy.ClientIP(),
		Method:   ctxCopy.Request.Method,
		FullPath: ctxCopy.FullPath(),
	}
	if kind == "Response" {
		accessLog.StatusCode = ctxCopy.Writer.Status()
	}
	accessLog.WriteLog()
}

func ginCustomLogger(c *gin.Context) {
	loggerTask("Request", c)
	c.Next()
	loggerTask("Response", c)
}

func ginLogHandler() gin.HandlerFunc {
	return ginCustomLogger
}

func (router *Router) Get() {
	router.R = gin.New()
	router.R.Use(ginLogHandler())
	gin.ForceConsoleColor()
}

func (router *Router) Run() {
	router.Get()
	router.InitRoutes()
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	if err := router.R.Run(":" + port); err != nil {
		log.Error(HeaderGin, "|", err)
		panic(err)
	}
}

func (router *Router) InitRoutes() {
	router.initHealthRoute()
	router.initApis()
}

func (router *Router) initHealthRoute() {
	router.R.GET("/", func(c *gin.Context) {
		//c.JSON(200, &utils.PostData{Kind: "Welcome"})
		c.SecureJSON(200, &utils.Body{Message: "Welcome"})
	})
	router.R.GET("/health", func(c *gin.Context) {
		c.SecureJSON(200, &utils.Body{Message: "pong"})
	})
}

func (router *Router) initApis() {
	api := router.R.Group("/api")
	{
		api.GET("/datetime", func(c *gin.Context) {
			response := utils.DatetimeApi()
			c.SecureJSON(200, response)
		})
		api.POST("/pray-schedule", func(c *gin.Context) {
			jsonBody := &utils.PrayScheduleReq{}
			if err := c.BindJSON(jsonBody); err != nil {
				log.Error(utils.HeaderPray, "|", err)
			} else {
				response := utils.GetSchedule(jsonBody)
				c.SecureJSON(200, &response)
			}
		})
	}
}
