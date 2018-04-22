package cmd

import (
	"github.com/spf13/cobra"
	"strings"
	"io/ioutil"
	"fmt"
	"github.com/pojozhang/sugar"
	"net/url"
	"playground-tool/util"
	"os"
)

type questionData struct {
	Data struct {
		Question struct {
			QuestionTitle     string
			TranslatedTitle   string
			Content           string
			TranslatedContent string
		}
	}
}

var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a problem markdown file",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		question := extractQuestion(url)
		file := question + ".md"

		if _, err := os.Stat(file); err == nil {
			return fmt.Errorf("%s already exists", file)
		}

		d := &questionData{}
		_, err := sugar.Get(questionUrl(question)).Read(d)
		if err != nil {
			return err
		}

		content := buildContent(url, d)
		ioutil.WriteFile(file, []byte(content), 0644)
		fmt.Printf("%s has been created.", file)
		return nil
	},
}

func questionUrl(title string) string {
	return "https://leetcode-cn.com/graphql?query=" + url.PathEscape(fmt.Sprintf(`query getQuestionDetail($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    questionTitle
    translatedTitle
    content
    translatedContent
  }
}
&operationName=getQuestionDetail&variables={"titleSlug":"%s"}`, title))
}

func extractQuestion(url string) string {
	const keyword = "problems/"
	index := strings.Index(url, keyword) + len(keyword)
	return url[index : strings.Index(url[index:], "/")+index]
}

func buildContent(url string, d *questionData) string {
	q := d.Data.Question
	return fmt.Sprintf("# [%s](%s)\n\n%s#### 实现\n", q.TranslatedTitle, url, util.ParseToMarkdown(q.TranslatedContent))
}
