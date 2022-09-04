package web

import (
	bingwallpaper "fatalisa-public-api/service/web/utils/bing-wallpaper"
	"fmt"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/subchen/go-log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var pageTitle = "Fatalisa Public API"

const webpagesDir = "service/web/pages"
const webBasePageDir = webpagesDir + "/base"

func Service(c *gin.Context, reqType string) {
	c.HTML(http.StatusOK, reqType, Index())
}

type BodyExample struct {
	Text1 string
	Text2 string
}

type FooterTexts struct {
	BgImgCopyright string
	BgImgUrl       string
	Text           string
	Year           string
}

func Index() *gin.H {
	bgImage := BackgroundImage()
	return &gin.H{
		"Title":    pageTitle,
		"Function": Example(),
		"BgImg":    bgImage.Url,
		"Body": BodyExample{
			Text1: "Welcome",
			Text2: "This is index page",
		},
		"Footer": FooterTexts{
			BgImgCopyright: bgImage.Copyright,
			BgImgUrl:       bgImage.CopyrightLink,
			Text:           pageTitle,
			Year:           fmt.Sprint(time.Now().Year()),
		},
	}
}

// LoadTemplates function referenced to gin multi-template documentation
// https://github.com/gin-contrib/multitemplate
func LoadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(webBasePageDir + "/*")
	if err != nil {
		log.Error("Error : ", err)
		return nil
	}

	var base []string
	base = append(base, webpagesDir+"/base/fullpage.gohtml")
	base = append(base, layouts...)

	pages, err := filepath.Glob(webpagesDir + "/*.gohtml")
	if err != nil {
		log.Error("Error: ", err)
	}
	for _, page := range pages {
		compiledPage := append(base, page)
		pageName := strings.Replace(filepath.Base(page), ".gohtml", "", -1)
		r.AddFromFiles(pageName, compiledPage...)
	}

	return r
}

func Example() string {
	return "Example"
}

func BackgroundImage() bingwallpaper.ImageData {
	return bingwallpaper.GetTodayWallpaper()
}
