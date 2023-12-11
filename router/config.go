package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	htmlEngine "github.com/gofiber/template/html/v2"
	"github.com/subchen/go-log"
	"os"
)

//const accessCounterKey = "accessCounter"

type Config struct {
	Fiber *fiber.App
}

func (router *Config) Get() {
	engine := htmlEngine.New("service/web/pages", ".gohtml")
	router.Fiber = fiber.New(
		fiber.Config{
			Views: engine,
		},
	)
	router.Fiber.Use(recover.New())
	router.Fiber.Use(fiberLogger.New(fiberLogger.Config{
		Format: "${time} [INFO] :" +
			fmt.Sprintf("%4s", "${status}") +
			"-" +
			fmt.Sprintf("%11s", "${latency}") +
			fmt.Sprintf("%7s", "${method}") +
			"${path}\n",
		TimeFormat: "2006/01/02 15:04:05 -0700",
	}))

	// 404 Handler
	//router.Fiber.Use(func(c *fiber.Ctx) error {
	//	return c.SendStatus(404) // => 404 "Not Found"
	//})
}

func (router *Config) Run() {
	router.Get()
	router.InitRoutes()
	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "80"
	}
	log.Info("Service running at port ", port)
	if err := router.Fiber.Listen(":" + port); err != nil {
		log.Error(err)
		panic(err)
	}
}

func (router *Config) InitRoutes() {
	router.initLandingRoute()
	router.initHealthRoute()
	router.initSwagger()
	router.versionChecker()
	router.initApis()
	router.initWebRoute()
	router.initServeFiles()
}

//func increaseAccessCounter() {
//	redis := config.InitRedis()
//	if len(redis.Client.Ping(context.Background()).Err().Error()) == 0 {
//		val := redis.Get(accessCounterKey)
//		currentValue, err := strconv.Atoi(val)
//		if err != nil {
//			log.Error(err)
//			currentValue = 0
//		}
//		currentValue++
//		redis.Set(accessCounterKey, strconv.Itoa(currentValue))
//	}
//}
