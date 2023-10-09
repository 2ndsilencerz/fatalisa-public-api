package model

import (
	"golang.org/x/net/html"
	"strings"
)

type ProvinceCode string

type ProvinceName string

type ProvinceMapping map[ProvinceName]ProvinceCode

func (province ProvinceMapping) ParseHtml(node *html.Node) {
	if node.Parent != nil && node.Parent.Data == "select" && len(node.Parent.Attr) > 0 &&
		node.Parent.Attr[0].Val == "search_prov" && node.Data == "option" {
		if !strings.HasPrefix(node.FirstChild.Data, "PILIH") {
			province[ProvinceName(node.FirstChild.Data)] = ProvinceCode(node.Attr[0].Val)
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		province.ParseHtml(c)
	}
}
