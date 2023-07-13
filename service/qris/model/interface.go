package model

type Request struct {
	Raw string `json:"raw" binding:"required"`
}

type General interface {
	Request
}
