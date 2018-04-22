package util

import (
	"strings"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func ParseToMarkdown(html string) string {
	html = strings.Replace(html, "\n\n", "\n", -1)
	md := ""
	d, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	parseToMarkdown(&md, d.Find("body").Nodes[0])
	return md
}

func parseToMarkdown(markdown *string, node *html.Node) {
	if node == nil {
		return
	}

	prefix, suffix := "", ""

	switch node.Type {
	case html.TextNode:
		if node.Parent.Data == "strong" && !strings.Contains(node.Data, "输入") &&
			!strings.Contains(node.Data, "输出") && !strings.Contains(node.Data, "解释") {
			*markdown += "**" + node.Data + "**"
		} else {
			*markdown += node.Data
		}
	case html.ElementNode:
		switch node.Data {
		case "code":
			prefix = "`"
			suffix = "`"
		case "pre":
			prefix = "\n```\n"
			suffix = "\n```\n"
		case "sup":
			prefix = "^"
		}
	}

	*markdown = *markdown + prefix

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		parseToMarkdown(markdown, c)
	}

	*markdown = *markdown + suffix
}
