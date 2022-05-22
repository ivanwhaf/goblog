package handlers

import (
	"github.com/gin-gonic/gin"
	"goblog/stores"
	"net/http"
	"strconv"
)

// TagsHandler /tags GET
func TagsHandler(c *gin.Context) {
	articles, _ := stores.ArticleStore.GetArticlesOrderByIdWithFields("id", "title", "tag")
	resMap := make(map[string][]map[string]string)
	for _, article := range articles {
		tag := article.Tag
		if _, ok := resMap[tag]; !ok {
			resMap[tag] = []map[string]string{}
		}
		m := map[string]string{
			"id":    strconv.FormatInt(article.Id, 10),
			"title": article.Title,
		}
		resMap[tag] = append(resMap[tag], m)
	}
	c.HTML(http.StatusOK, "tags.html", gin.H{"resMap": resMap})
}
