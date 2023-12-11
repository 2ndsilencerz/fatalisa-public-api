package web

import (
	"bytes"
	"fatalisa-public-api/service/web/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/subchen/go-log"
	"html/template"
)

var pageTitle = "Fatalisa Public API"

func loadTemplate(c *fiber.Ctx, basePage bool) *template.Template {
	var templateFile *template.Template
	var err error
	if c != nil && !basePage {
		templateFile, err = template.ParseFiles("./service/web/pages/" + c.Params("path", "index") + ".gohtml")
	} else {
		templateFile, err = template.ParseFiles(
			"./service/web/pages/base/footer.gohtml",
			"./service/web/pages/base/navbar.gohtml",
			"./service/web/pages/base/fullpage.gohtml",
		)
	}
	utils.ErrorHandler(err)
	return templateFile
}

func Layout() string {
	templateRaw := loadTemplate(nil, true)
	//templateNavbar := loadTemplate()
	//templateFooter :=
	var buf bytes.Buffer
	err := templateRaw.Execute(&buf, "layout")
	log.Info("Executing template")
	utils.ErrorHandler(err)
	return buf.String()
}

func MainLayout(c *fiber.Ctx) string {
	templateRaw := loadTemplate(c, false)
	var buf bytes.Buffer
	err := templateRaw.Execute(&buf, "layout")
	log.Info("Executing template main")
	utils.ErrorHandler(err)
	return buf.String()
}

func Index() *fiber.Map {
	body, footer := Page()
	contentMap := fiber.Map{
		"Title":  pageTitle,
		"Footer": footer,
		"BgImg":  footer.BgImgUrl,
		"Body":   body,
	}

	return &contentMap
}

//func DynamicContent(c *fiber.Ctx) *fiber.Map {
//	template := loadTemplate(c)
//	template, err = template.Parse()
//}
