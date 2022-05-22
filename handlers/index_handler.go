package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/core"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"strconv"
	"time"
)

// IndexHandler / GET
func IndexHandler(c *gin.Context) {
	strPage := c.DefaultQuery("page", "1")
	page, _ := strconv.ParseInt(strPage, 10, 64)
	articlesCount, _ := stores.ArticleStore.GetArticlesCount()
	visitorsCount, _ := stores.VisitorStore.GetVisitorsCount()

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
		go AddVisitorRecord(&core.Visitor{
			Ip:        ip,
			VisitDate: visitDate,
			Platform:  env.Platform,
			Browser:   env.Browser,
			Version:   env.Version,
		})
	}

	var start int64
	if page > 1 {
		start = (page - 1) * 10
	} else {
		start = 0
	}

	var popularArticles []*core.Article
	popularArticles, _ = stores.ArticleStore.GetMostPopularArticles(5)
	var articles []*core.Article
	articles, _ = stores.ArticleStore.GetArticlesOrderByIdOffsetLimit(start, 10)

	for _, article := range articles {
		if len(article.ContentText) > 300 {
			article.ContentText = article.ContentText[:300]
		}
	}

	session := sessions.Default(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"username":        session.Get("username"),
		"nickname":        session.Get("nickname"),
		"page":            page,
		"pageSub1":        page - 1,
		"pagePlu1":        page + 1,
		"pageMul10":       page * 10,
		"articlesCount":   articlesCount,
		"visitorsCount":   visitorsCount + 1,
		"articles":        articles,
		"popularArticles": popularArticles,
	})
}

// RobotsHandler /robots.txt GET
func RobotsHandler(c *gin.Context) {
	c.File(config.GetConfig().File.RobotsTxtPath)
}

func GetLocation(ip string) string {
	var location string
	ipLocation, _ := stores.IpLocationStore.GetIpLocationByIp(ip)
	if ipLocation.Location == "" {
		location = util.CrawlIpLocation(ip)
	} else {
		location = ipLocation.Location
	}
	return location
}

func AddVisitorRecord(v *core.Visitor) {
	v.Location = GetLocation(v.Ip)
	_ = stores.VisitorStore.AddVisitor(v)
}
