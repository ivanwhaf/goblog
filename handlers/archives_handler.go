package handlers

import (
	"github.com/gin-gonic/gin"
	"goblog/stores"
	"net/http"
	"strconv"
	"strings"
)

// ArchivesHandler /archives GET
func ArchivesHandler(c *gin.Context) {
	articles, _ := stores.ArticleStore.GetArticlesOrderById()
	resMap := make(map[string][]map[string]string)
	for _, article := range articles {
		createDate := article.CreateDate
		ymd := createDate.String()
		ymd = strings.Split(ymd, " ")[0]
		year := ymd[0:4]
		md := ymd[5:10]
		if articlesSlc, ok := resMap[year]; ok {
			articlesSlc = append(articlesSlc, map[string]string{
				"Id":    strconv.FormatInt(article.Id, 10),
				"Title": article.Title,
				"Tag":   article.Tag,
				"Md":    md,
			})
			resMap[year] = articlesSlc
		} else {
			articlesSlc = []map[string]string{}
			articlesSlc = append(articlesSlc, map[string]string{
				"Id":    strconv.FormatInt(article.Id, 10),
				"Title": article.Title,
				"Tag":   article.Tag,
				"Md":    md,
			})
			resMap[year] = articlesSlc
		}
	}

	c.HTML(http.StatusOK, "archives.html", gin.H{
		"resMap": resMap,
	})
}
