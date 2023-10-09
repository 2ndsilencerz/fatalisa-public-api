package model

import (
	"encoding/json"
	"github.com/subchen/go-log"
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
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal(raw, &schedule)
	if err != nil {
		log.Error(err)
		return
	}
}
