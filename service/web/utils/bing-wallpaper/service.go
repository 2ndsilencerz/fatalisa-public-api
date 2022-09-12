package bing_wallpaper

import (
	"encoding/json"
	"github.com/subchen/go-log"
	"io"
	"net/http"
)

const (
	BingBaseUrl     = "https://www.bing.com"
	bingWallBaseURL = BingBaseUrl + "/HPImageArchive.aspx?format=js&idx=0&n=1"
	errorMsg        = "Error while getting wallpaper"
)

type jsonAPI struct {
	Images []ImageData `json:"images"`
}

type ImageData struct {
	Url           string `json:"url"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightlink"`
}

func errorExist(err error) bool {
	if err != nil {
		log.Error(errorMsg, err)
		return true
	}
	return false
}

func GetTodayWallpaper() ImageData {
	var imageData ImageData
	apiResponse, err := http.Get(bingWallBaseURL)
	if errorExist(err) {
		return imageData
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(apiResponse.Body)
	rawJson, err := io.ReadAll(apiResponse.Body)

	data := &jsonAPI{}
	err = json.Unmarshal(rawJson, data)
	if errorExist(err) {
		return imageData
	}
	imageData = data.Images[0]
	imageData.Url = BingBaseUrl + imageData.Url
	return imageData
}
