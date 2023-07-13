package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	htmlEngine "github.com/gofiber/template/html/v2"
	"github.com/subchen/go-log"
	"html/template"
	"os"
)

//const accessCounterKey = "accessCounter"

type Config struct {
	Fiber *fiber.App
}

func (router *Config) Get() {
	engine := htmlEngine.New("./service/web/pages", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	router.Fiber = fiber.New(
		fiber.Config{
			Views: engine,
		},
	)
	router.Fiber.Use(recover.New())
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
	//router.initLandingRoute()
	router.initHealthRoute()
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
