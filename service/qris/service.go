package qris

import (
	"fatalisa-public-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
)

func ParseMpmService(c *gin.Context) *MpmData {
	req := MpmRequest{}
	if len(c.Param("raw")) > 0 {
		req.Raw = c.Param("raw")
	} else if err := c.BindJSON(&req); err != nil {
		log.Error(err)
	} else {
		log.Info(utils.Jsonify(req))
	}

	res := MpmData{}
	res.GetData(req.Raw)
	log.Info(utils.Jsonify(res))
	go res.SaveToDB()
	return &res
}

func ParseCpmService(c *gin.Context) *CpmData {
	req := CpmRequest{}
	if err := c.BindJSON(&req); err != nil {
		log.Error(err)
	} else {
		log.Info(req)
	}

	res := CpmData{}
	res.GetData(req.Raw)
	log.Info(utils.Jsonify(res))
	go res.SaveToDB()
	return &res
}
