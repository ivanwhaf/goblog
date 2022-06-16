package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/core"
	"goblog/services"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// IndexHandler / GET
func IndexHandler(c *gin.Context) {
	session := sessions.Default(c)

	strPage := c.DefaultQuery("page", "1")
	page, _ := strconv.ParseInt(strPage, 10, 64)
	articlesCount, _ := stores.ArticleStore.GetArticlesCount()

	if (page-1)*10 >= articlesCount || page < 1 {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "404",
		})
		return
	}

	if page == 1 {
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		env := util.ParseUserAgent(ua)
		visitDate := time.Now()
		go services.AddVisitorRecord(&core.Visitor{
			Ip:        ip,
			VisitDate: visitDate,
			Platform:  env.Platform,
			Browser:   env.Browser,
			Version:   env.Version,
		})
	}

	var wg sync.WaitGroup
	var articlesSlc []map[string]string
	var popularArticles []*core.Article
	var visitorsCount int64

	wg.Add(3)
	go func() {
		var start int64
		if page > 1 {
			start = (page - 1) * 10
		} else {
			start = 0
		}
		articles, _ := stores.ArticleStore.GetArticlesOrderByIdOffsetLimit(start, 10)
		articlesSlc = services.GetArticlesSlc(articles)
		wg.Done()
	}()

	go func() {
		popularArticles, _ = stores.ArticleStore.GetMostPopularArticles(5)
		wg.Done()
	}()

	go func() {
		// visitorsCount, _ = stores.VisitorStore.GetVisitorsCount()
		visitorsCount = stores.RedisStore.GetVisitorsCount()
		stores.RedisStore.IncreaseVisitorsCount()
		wg.Done()
	}()

	wg.Wait()

	c.HTML(http.StatusOK, "index.html", gin.H{
		"username":        session.Get("username"),
		"nickname":        session.Get("nickname"),
		"page":            page,
		"pageSub1":        page - 1,
		"pagePlu1":        page + 1,
		"pageMul10":       page * 10,
		"articlesCount":   articlesCount,
		"visitorsCount":   visitorsCount + 1,
		"articles":        articlesSlc,
		"popularArticles": popularArticles,
	})
}

// RobotsHandler /robots.txt GET
func RobotsHandler(c *gin.Context) {
	c.File(config.GetConfig().File.RobotsTxtPath)
}
