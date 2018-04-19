package cmd

import (
	"github.com/spf13/cobra"
	"net/http"
	"strings"
	"io/ioutil"
	"fmt"
	"html"
	"github.com/berryland/sugar"
)

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a problem markdown file",
	RunE: func(cmd *cobra.Command, args []string) error {

		type data struct {
		}
		url := "https://leetcode-cn.com/problems/string-to-integer-atoi/description/"
		question := extractQuestion(url)

		d := &data{}
		sugar.Get(questionUrl(question)).Read(d)
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		title := extractTitle(doc.Find("title").Text())
		problem, _ := doc.Find("meta[name='description']").Attr("content")
		content := buildContent(title, url, problem)
		ioutil.WriteFile(fileName, []byte(content), 0644)
		return nil
	},
}

func questionUrl(title string) string {
	return fmt.Sprintf(`https://leetcode-cn.com/graphql?query=query getQuestionDetail($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionTitle
    translatedTitle
    content
    translatedContent
  }
}
&operationName=getQuestionDetail&variables={"titleSlug":"%s"}`, title)
}

func extractQuestion(url string) string {
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
