package pray_schedule

import (
	"fatalisa-public-api/service/pray-schedule/kemenag"
	"fatalisa-public-api/service/pray-schedule/model"
	"fatalisa-public-api/service/pray-schedule/pkpu"
	"fatalisa-public-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/subchen/go-log"
	"os"
	"time"
)

const (
	ProviderPkpu    = "PKPU"
	ProviderKemenag = "Kemenag"
)

var provider string

func init() {
	provider, _ = os.LookupEnv("PROVIDER")
}

// GetSchedule used to get pray schedule of requested city and date if provided
func GetSchedule(c *fiber.Ctx) *model.Response {
	req := model.Request{}
	if err := c.BodyParser(&req); err != nil {
		req.City = c.Params("city")
		req.Date = time.Now().Format("2006/01/02")
	}
	log.Info(utils.Jsonify(req))

	var res *model.Response
	if provider == ProviderPkpu {
		res = pkpu.GetData(&req)
	} else {
		res = kemenag.GetData(&req)
	}
	log.Info(utils.Jsonify(res))
	return res
}

func GetCityList() interface{} {
	if provider == ProviderPkpu {
		return pkpu.GetCityList()
	}
	return kemenag.GetCityList()
}
