package cpm

type Request struct {
	Raw string `json:"raw" binding:"required"`
}
