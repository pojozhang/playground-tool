package cmd

import (
	"github.com/spf13/cobra"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"html"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a problem markdown file",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://leetcode-cn.com/problems/string-to-integer-atoi/description/"
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return err
		}

		title := extractTitle(doc.Find("title").Text())
		problem, _ := doc.Find("meta[name='description']").Attr("content")
		content := buildContent(title, url, problem)
		fileName := extractFileName(url)
		ioutil.WriteFile(fileName, []byte(content), 0644)
		return nil
	},
}

func extractFileName(url string) string {
	const keyword = "problems/"
	index := strings.Index(url, keyword) + len(keyword)
	return url[index:strings.Index(url[index:], "/")+index] + ".md"
}

func extractTitle(title string) string {
	return strings.TrimSpace(title[:strings.LastIndex(title, "-")])
}

func buildContent(title, url, problem string) string {
	p := html.UnescapeString(problem)
	return fmt.Sprintf("# [%s](%s)\n\n %s\n #### 实现\n", title, url, p)
}
