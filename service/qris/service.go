package qris

import (
	qrisModel "fatalisa-public-api/service/qris/model"
	"fatalisa-public-api/service/qris/model/cpm"
	"fatalisa-public-api/service/qris/model/mpm"
	"fatalisa-public-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/subchen/go-log"
)

// ParseMpm used to parse MPM data
func ParseMpm(c *fiber.Ctx) *mpm.Data {
	req := qrisModel.Request{}
	if len(c.Params("raw")) > 0 {
		req.Raw = c.Params("raw")
	} else if err := c.BodyParser(&req); err != nil {
		log.Error(err)
		return nil
	} else {
		log.Info(utils.Jsonify(req))
	}

	res := mpm.Data{}
	res.Parse(req.Raw)
	log.Info(utils.Jsonify(res))
	go res.SaveToDB()
	return &res
}

// ParseCpm used to parse CPM data
func ParseCpm(c *fiber.Ctx) *cpm.Data {
	req := qrisModel.Request{}
	if err := c.BodyParser(&req); err != nil {
		log.Error(err)
		return nil
	} else {
		log.Info(req)
	}

	res := cpm.Data{}
	res.Parse(req.Raw)
	log.Info(utils.Jsonify(res))
	go res.SaveToDB()
	return &res
}
