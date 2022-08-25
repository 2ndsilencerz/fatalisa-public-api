package mpm

type Request struct {
	Raw string `json:"raw" binding:"required"`
}
