package services

import (
	"goblog/core"
	"strconv"
)

func GetArticlesSlc(articles []*core.Article) []map[string]string {
	var articlesSlc []map[string]string
	for _, article := range articles {
		runeText := []rune(article.ContentText)
		if len(runeText) > 300 {
			article.ContentText = string(runeText[:300])
		}
		m := make(map[string]string)
		m["Id"] = strconv.FormatInt(article.Id, 10)
		m["Tag"] = article.Tag
		m["Title"] = article.Title
		m["Subtitle"] = article.Subtitle
		m["ContentText"] = article.ContentText
		m["ReadCount"] = strconv.FormatInt(article.ReadCount, 10)
		m["Author"] = article.Author
		m["CreateDate"] = article.CreateDate.Format("2006-01-02 15:03:04")
		articlesSlc = append(articlesSlc, m)
	}
	return articlesSlc
}
