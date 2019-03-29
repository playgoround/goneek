package parser

import (
	"bytes"
	"golang.org/x/net/html"
)

// separated indexed array
type sepIdxArr struct {
	index int
	array []string
}

func getAttrFromNode(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getTextFromNode(node *html.Node) string {
	var b bytes.Buffer
	html.Render(&b, node)
	return b.String()
}
