package web

import (
	"github.com/subchen/go-log"
	"html/template"
)

//var pageTitle = "Fatalisa Public API"

func Index() *template.Template {
	//body := BodyExample{
	//	Text1: "Welcome",
	//	Text2: "This is index page",
	//}
	tmpl, _ := template.New("index").ParseFiles("service/web/template/index.html")
	log.Debug("tmpl: ", tmpl)
	return tmpl
}
