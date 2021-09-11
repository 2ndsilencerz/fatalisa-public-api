package config

import (
	"encoding/json"
	"fatalisa-public-api/database"
	"fatalisa-public-api/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pieterclaerhout/go-log"
	"os"
	"time"
)

type Router struct {
	R *gin.Engine
}

func (router *Router) get() {
	router.R = gin.New()
	router.R.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	gin.ForceConsoleColor()
}

func (router *Router) Run() {
	router.get()
	router.initRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	err := router.R.Run(":" + port)
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func (router *Router) initRoutes() {
	router.initHealthRoute()
	router.initApis()
}

func (router *Router) initHealthRoute() {
	router.R.GET("/", func(c *gin.Context) {
		//c.JSON(200, &utils.Body{Message: "Welcome"})
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
			tmp, err := json.Marshal(response)
			if err != nil {
				log.Error(err)
				errorLog := database.ErrorLog{}
				errorLog.Write(err)
			}
			accessLog := &database.AccessLog{
				IP:           c.ClientIP(),
				Method:       "",
				FullPath:     c.FullPath(),
				ResponseCode: 200,
				Response:     string(tmp),
			}
			accessLog.WriteLog()
		})
	}
}
