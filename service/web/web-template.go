package web

import bingwallpaper "fatalisa-public-api/service/web/utils/bing-wallpaper"

const (
// webpagesDir    = "./service/web/pages"
// webBasePageDir = webpagesDir + "/base"
)

type BodyExample struct {
	Title    string
	Text1    string
	Text2    string
	Function interface{}
}

type FooterTexts struct {
	BgImgCopyright string
	BgImgUrl       string
	Text           string
	Year           string
	Version        string
}

func BackgroundImage() bingwallpaper.ImageData {
	return bingwallpaper.GetTodayWallpaper()
}

func WebTemplate() (*BodyExample, *FooterTexts) {
	footer := FooterTexts{
		BgImgUrl:       BackgroundImage().Url,
		BgImgCopyright: BackgroundImage().Copyright,
		Text:           "Fatalisa Public API",
		Year:           "2023",
		Version:        "1.0.0",
	}

	body := BodyExample{
		Title:    "Fatalisa Public API",
		Text1:    "Hello World",
		Text2:    "Hello World 2",
		Function: "Hello World Map",
	}

	return &body, &footer
}
