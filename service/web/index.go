package web

import (
	"github.com/gin-gonic/gin"
)

var pageTitle = "Fatalisa Public API"

func Index() *gin.H {
	body := BodyExample{
		Text1: "Welcome",
		Text2: "This is index page",
	}

	return Template(body)
}
