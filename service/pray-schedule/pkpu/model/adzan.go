package model

import (
	"fatalisa-public-api/service/pray-schedule/model"
)

// Adzan used as model to fetch schedule data from xml
type Adzan struct {
	Version   string `xml:"version"`
	Site      string `xml:"site"`
	Country   string `xml:"country"`
	City      string `xml:"city"`
	Parameter `xml:"parameter"`
	Data      []model.Response `xml:"data" json:"data"`
}
