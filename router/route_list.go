package router

import (
	_ "fatalisa-public-api/docs"
	"fatalisa-public-api/service/common"
	praySchedule "fatalisa-public-api/service/pray-schedule"
	"fatalisa-public-api/service/qris"
	"fatalisa-public-api/service/web"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (router *Config) initLandingRoute() {
	router.Fiber.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", web.Index())
	})
}

// initHealthRoute godoc
//
//	@Summary		HealthCheck
//	@Description	Health Check
//	@Tags			Default
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	common.Body
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (router *Config) initHealthRoute() {
	router.Fiber.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(&common.Body{Message: "pong"})
	})
}

// versionChecker godoc
//
//	@Summary		VersionChecker
//	@Description	Show Version
//	@Tags			Default
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	common.Body
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/version [get]
func (router *Config) versionChecker() {
	router.Fiber.Get("/version", func(c *fiber.Ctx) error {
		return c.JSON(common.VersionChecker())
	})
}

// PrayScheduleCityList godoc
//
//	@Summary		PrayScheduleCityList
//	@Description	Get Available City List
//	@Tags			Pray-Schedule
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.CityList
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/pray-schedule/city-list [get]
func (router *Config) PrayScheduleCityList() {
	router.Fiber.Group("/api").Get("/pray-schedule/city-list", func(c *fiber.Ctx) error {
		return c.JSON(praySchedule.GetCityList())
	})
}

// PrayScheduleCity godoc
//
//	@Summary		PrayScheduleCity
//	@Description	Get Schedule By City
//	@Tags			Pray-Schedule
//	@Accept			json
//	@Produce		json
//	@Param			city	path		string	true	"city"
//	@Success		200	{object}	model.Response
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/pray-schedule/:city [get]
func (router *Config) PrayScheduleCity() {
	router.Fiber.Group("/api").Get("/pray-schedule/:city", func(c *fiber.Ctx) error {
		return c.JSON(praySchedule.GetSchedule(c))
	})
}

// PrayScheduleCityPost godoc
//
//	@Summary		PrayScheduleCityPost
//	@Description	Get Schedule By City and Date
//	@Tags			Pray-Schedule
//	@Accept			json
//	@Produce		json
//	@Param			city	body	interface{}	true	"data"
//	@Success		200	{object}	model.Response
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/pray-schedule [post]
func (router *Config) PrayScheduleCityPost() {
	router.Fiber.Group("/api").Post("/pray-schedule", func(c *fiber.Ctx) error {
		return c.JSON(praySchedule.GetSchedule(c))
	})
}

// ParseMpmGet godoc
//
//	@Summary		ParseMpm
//	@Description	Parse MPM by parameter
//	@Tags			MPM
//	@Accept			json
//	@Produce		json
//	@Param			raw	path		string	true	"raw"
//	@Success		200	{object}	mpm.Data
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/qris/mpm/:raw [get]
func (router *Config) ParseMpmGet() {
	router.Fiber.Group("/api/qris").Get("/mpm/:raw", func(c *fiber.Ctx) error {
		return c.JSON(qris.ParseMpm(c))
	})
}

// ParseMpmPost godoc
//
//	@Summary		ParseMpm
//	@Description	Parse MPM by raw body
//	@Tags			MPM
//	@Accept			json
//	@Produce		json
//	@Param			raw	body 		mpm.Request true	"data"
//	@Success		200	{object}	mpm.Data
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/qris/mpm [post]
func (router *Config) ParseMpmPost() {
	router.Fiber.Group("/api/qris").Post("/mpm", func(c *fiber.Ctx) error {
		return c.JSON(qris.ParseMpm(c))
	})
}

// ParseCpmPost godoc
//
//	@Summary		ParseCpm
//	@Description	Parse CPM by raw body
//	@Tags			CPM
//	@Accept			json
//	@Produce		json
//	@Param			raw	body		cpm.Request true	"data"
//	@Success		200	{object}	cpm.Data
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/api/qris/cpm [post]
func (router *Config) ParseCpmPost() {
	router.Fiber.Group("/api/qris").Post("/cpm", func(c *fiber.Ctx) error {
		return c.JSON(qris.ParseCpm(c))
	})
}

func (router *Config) initApis() {
	router.PrayScheduleCityList()
	router.PrayScheduleCity()
	router.PrayScheduleCityPost()
	router.ParseMpmGet()
	router.ParseMpmPost()
	router.ParseCpmPost()
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

func (router *Config) initSwagger() {
	router.Fiber.Get("/swagger/*", swagger.HandlerDefault) // default
}
