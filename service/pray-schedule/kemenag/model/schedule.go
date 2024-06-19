package model

import (
	"encoding/json"
	utils2 "fatalisa-public-api/service/web/utils"
	"io"
	"net/http"
)

type Schedule struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Prov    string  `json:"prov"`
	Kabko   string  `json:"kabko"`
	Data    Dailies `json:"data"`
}

func (schedule *Schedule) Parse(response *http.Response) {
	// parse response to Schedule
	raw, err := io.ReadAll(response.Body)
	if err, _ := utils2.ErrorHandler(err); err {
		return
	}
	err = json.Unmarshal(raw, &schedule)
	if err, _ := utils2.ErrorHandler(err); err {
		return
	}
}
