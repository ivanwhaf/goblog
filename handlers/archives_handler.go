package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/stores"
	"net/http"
	"strconv"
	"strings"
)

// ArchivesHandler /archives GET
func ArchivesHandler(c *gin.Context) {
	session := sessions.Default(c)
	articles, _ := stores.ArticleStore.GetArticlesOrderById()
	resMap := make(map[string][]map[string]string)
	for _, article := range articles {
		createDate := article.CreateDate.Format("2006-01-02 15:03:04")
		slc := strings.Split(createDate, " ")
		ymd := slc[0]
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
		"resMap":   resMap,
		"username": session.Get("username"),
		"nickname": session.Get("nickname"),
	})
}
