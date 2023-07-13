package router

import (
	"fatalisa-public-api/service/common"
	praySchedule "fatalisa-public-api/service/pray-schedule"
	"fatalisa-public-api/service/qris"
	"fatalisa-public-api/service/web"
	"github.com/gofiber/fiber/v2"
)

func (router *Config) initLandingRoute() {
	router.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", web.Index())
	})
}

func (router *Config) initHealthRoute() {
	router.Fiber.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(&common.Body{Message: "pong"})
	})
}

func (router *Config) versionChecker() {
	router.Fiber.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(common.VersionChecker())
	})
}

func (router *Config) initApis() {
	api := router.Fiber.Group("/api")
	{
		api.Get("/datetime", func(c *fiber.Ctx) error {
			return c.JSON(common.DateTimeApi())
		})
		api.Get("/pray-schedule/city-list", func(c *fiber.Ctx) error {
			return c.JSON(praySchedule.GetCityList())
		})
		api.Get("/pray-schedule/:city", func(c *fiber.Ctx) error {
			return c.JSON(praySchedule.GetSchedule(c))
		})
		api.Post("/pray-schedule", func(c *fiber.Ctx) error {
			return c.JSON(praySchedule.GetSchedule(c))
		})
		qrisGroup := api.Group("/qris")
		{
			qrisGroup.Get("/mpm/:raw", func(c *fiber.Ctx) error {
				return c.JSON(qris.ParseMpm(c))
			})
			qrisGroup.Post("/mpm", func(c *fiber.Ctx) error {
				return c.JSON(qris.ParseMpm(c))
			})
			qrisGroup.Post("/cpm", func(c *fiber.Ctx) error {
				return c.JSON(qris.ParseCpm(c))
			})
		}
	}
}

func (router *Config) initWebRoute() {
	api := router.Fiber.Group("/web")
	{
		api.Get("/pray-schedule", func(c *fiber.Ctx) error {
			return c.JSON("none")
		})
	}
}

func (router *Config) initServeFiles() {
	router.Fiber.Static("/", "./service/web/pages")
	router.Fiber.Static("/css", "./service/web/pages/css")
	router.Fiber.Static("/img", "./service/web/pages/img")
}
