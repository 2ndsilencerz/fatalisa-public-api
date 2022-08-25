package qris

import (
	"fatalisa-public-api/service/qris/model/cpm"
	"fatalisa-public-api/service/qris/model/mpm"
	"fatalisa-public-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
)

// ParseMpm used to parse MPM data
func ParseMpm(c *gin.Context) *mpm.Data {
	req := mpm.Request{}
	if len(c.Param("raw")) > 0 {
		req.Raw = c.Param("raw")
	} else if err := c.BindJSON(&req); err != nil {
		log.Error(err)
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
func ParseCpm(c *gin.Context) *cpm.Data {
	req := cpm.Request{}
	if err := c.BindJSON(&req); err != nil {
		log.Error(err)
	} else {
		log.Info(req)
	}

	res := cpm.Data{}
	res.Parse(req.Raw)
	log.Info(utils.Jsonify(res))
	go res.SaveToDB()
	return &res
}
