package model

type Request struct {
	City string `json:"city" binding:"required"`
	Date string `json:"date" binding:"required"`
}

type City struct {
	Name string `json:"cityName"`
}
