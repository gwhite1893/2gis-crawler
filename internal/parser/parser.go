package parser

import (
	"io"
	"log"

	"golang.org/x/net/html"
)

// GetTagValue возвращает значение тега title страницы
func GetTagValue(r io.Reader, tagName string) string {
	doc, err := html.Parse(r)
	if err != nil {
		log.Println(err)

		return ""
	}

	node := getNode(doc, tagName)
	if node == nil || node.FirstChild == nil {
		return ""
	}

	return node.FirstChild.Data
}

func getNode(doc *html.Node, tagName string) *html.Node {
	var (
		res    *html.Node
		finder func(*html.Node)
	)

	finder = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == tagName {
			res = node

			return
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}

	finder(doc)

	return res
}
