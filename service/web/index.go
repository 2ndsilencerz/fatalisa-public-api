package web

import (
	"github.com/gofiber/fiber/v2"
)

var pageTitle = "Fatalisa Public API"

func Index() *fiber.Map {
	body, footer := WebTemplate()
	contentMap := fiber.Map{
		"Title":  pageTitle,
		"Footer": footer,
		"BgImg":  footer.BgImgUrl,
		"Body":   body,
	}

	return &contentMap
}
