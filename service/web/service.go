package web

import (
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
	"net/http"
	"path/filepath"
	"time"
)

var pageTitle = "Fatalisa Public API"

const webpagesDir = "service/web/pages"
const webBasePageDir = webpagesDir + "/base"

func Service(c *gin.Context, reqType string) {
	c.HTML(http.StatusOK, reqType, Index())
}

func Index() *gin.H {
	return &gin.H{
		"Title": pageTitle,
		"Body": struct {
			Text1 string
			Text2 string
		}{
			Text1: "Welcome",
			Text2: "This is index page",
		},
		"Footer": struct {
			Text string
			Year string
		}{
			Text: pageTitle,
			Year: fmt.Sprint(time.Now().Year()),
		},
	}
}

// LoadTemplates function referenced to gin multi-template documentation
// https://github.com/gin-contrib/multitemplate
func LoadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(webBasePageDir + "/*")
	if err != nil {
		log.Error(err)
		return nil
	}

	var index []string
	index = append(index, webpagesDir+"/index.gohtml")
	index = append(index, layouts...)
	r.AddFromFiles("index", index...)

	return r
}
