package model

import "golang.org/x/net/html"

type CityCode string

type CityName string

type CityMapping map[CityName]CityCode

func (city CityMapping) ParseHtml(node *html.Node) {
	if node.Data == "option" {
		city[CityName(node.FirstChild.Data)] = CityCode(node.Attr[0].Val)
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		city.ParseHtml(c)
	}
}
