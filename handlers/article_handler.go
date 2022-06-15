package handlers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pochard/commons/randstr"
	"goblog/core"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ArticleDetailHandler /article/:tag/:id GET
func ArticleDetailHandler(c *gin.Context) {
	session := sessions.Default(c)
	articleIdStr := c.Param("id")
	tag := c.Param("tag")
	articleId := util.StringToInt64(articleIdStr)
	article, _ := stores.ArticleStore.GetArticleById(articleId)
	if article.Tag != tag {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
		return
	}

	reader := strings.NewReader(article.ContentHtml)
	doc, _ := goquery.NewDocumentFromReader(reader)
	type Node struct {
		Children []*Node
		Text     string
		Id       string
	}
	var tagNodeLst []*Node
	var lastH1Node *Node
	doc.Find("body").Children().Each(func(i int, selection *goquery.Selection) {
		if selection.Nodes[0].Data == "h1" {
			node := &Node{}
			node.Id, _ = selection.Attr("id")
			node.Text = selection.Text()
			tagNodeLst = append(tagNodeLst, node)
			lastH1Node = node
		} else if selection.Nodes[0].Data == "h2" {
			node := &Node{}
			node.Id, _ = selection.Attr("id")
			node.Text = selection.Text()
			if lastH1Node != nil {
				lastH1Node.Children = append(lastH1Node.Children, node)
			} else {
				tagNodeLst = append(tagNodeLst, node)
			}
		}
	})

	articles, _ := stores.ArticleStore.GetArticlesOrderByIdWithFields("id", "title", "tag")

	var previousId, nextId string
	var previousTag, nextTag string
	var previousTagTitle, nextTagTitle string
	for idx, article := range articles {
		if article.Id == articleId {
			if idx > 0 {
				previousArticle := articles[idx-1]
				previousId = strconv.FormatInt(previousArticle.Id, 10)
				previousTag = previousArticle.Tag
				previousTagTitle = "[" + previousTag + "]" + previousArticle.Title
			}
			if idx < len(articles)-1 {
				nextArticle := articles[idx+1]
				nextId = strconv.FormatInt(nextArticle.Id, 10)
				nextTag = nextArticle.Tag
				nextTagTitle = "[" + nextTag + "]" + nextArticle.Title
			}
		}
	}

	comments, _ := stores.CommentStore.GetCommentsByArticleId(articleId)
	var commentsSlc []map[string]string
	for _, comment := range comments {
		commentMap := map[string]string{
			"replyName":   comment.ReplyName,
			"content":     comment.Content,
			"commentDate": comment.CommentDate.String(),
			"platform":    comment.Platform,
			"browser":     comment.Browser,
			"ip":          comment.Ip,
			"location":    comment.Location,
			"reviewerId":  strconv.FormatInt(comment.ReviewerId, 10),
		}
		if comment.ReviewerId == -1 {
			commentMap["avatar"] = "/files/avatar/head.jfif?r=" + randstr.RandomAlphabetic(2)
		} else {
			admin, _ := stores.AdminStore.GetAdminById(comment.ReviewerId)
			commentMap["avatar"] = "/files/avatar/" + admin.Avatar + "?r=" + randstr.RandomAlphabetic(2)
		}
		commentsSlc = append(commentsSlc, commentMap)
	}

	_ = stores.ArticleStore.UpdateArticleReadCount(articleId)

	c.HTML(http.StatusOK, "article.html", gin.H{
		"articleId":        article.Id,
		"title":            article.Title,
		"subtitle":         article.Subtitle,
		"tag":              article.Tag,
		"contentMd":        article.ContentMd,
		"contentHtml":      article.ContentHtml,
		"contentText":      article.ContentText,
		"author":           article.Author,
		"authorId":         article.AuthorId,
		"readCount":        article.ReadCount,
		"createDate":       article.CreateDate.Format("2006-01-02 15:03:04"),
		"lastEditDate":     article.LastEditDate.Format("2006-01-02 15:03:04"),
		"comments":         commentsSlc,
		"previousId":       previousId,
		"previousTag":      previousTag,
		"previousTagTitle": previousTagTitle,
		"nextId":           nextId,
		"nextTag":          nextTag,
		"nextTagTitle":     nextTagTitle,
		"tagNodeLst":       tagNodeLst,
		"username":         session.Get("username"),
		"nickname":         session.Get("nickname"),
	})
}

// ArticleAddHandler /article/add GET
func ArticleAddHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	} else if !util.ElementInSlice(authority, []int8{1, 2}) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	tags, _ := stores.ArticleStore.GetAllArticlesTags()
	c.HTML(http.StatusOK, "add_article.html", gin.H{"tags": tags})
}

// ArticleEditHandler /article/:tag/:id/edit
func ArticleEditHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	} else if !util.ElementInSlice(authority, []int8{1, 2}) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	id := c.Param("id")
	tag := c.Param("tag")

	article, _ := stores.ArticleStore.GetArticleById(util.StringToInt64(id))
	if article.Tag != tag {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	tags, _ := stores.ArticleStore.GetAllArticlesTags()

	c.HTML(http.StatusOK, "edit_article.html", gin.H{
		"article": article,
		"tags":    tags,
	})
}

// ArticleSearchHandler /article/search GET
func ArticleSearchHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", gin.H{})
}

// ApiArticleAddHandler /api/article POST
func ApiArticleAddHandler(c *gin.Context) {
	session := sessions.Default(c)
	userid := session.Get("userid")
	username := session.Get("username")
	authority := session.Get("authority")
	if userid == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.String(200, "-1")
		return
	}
	title := c.DefaultPostForm("title", "")
	subtitle := c.DefaultPostForm("subtitle", "")
	tag := c.DefaultPostForm("tag", "")
	contentMd := c.DefaultPostForm("content_md", "")
	contentHtml := c.DefaultPostForm("content_html", "")
	reader := strings.NewReader(contentHtml)
	doc, _ := goquery.NewDocumentFromReader(reader)
	contentText := doc.Text()

	err := stores.ArticleStore.AddArticle(&core.Article{
		Title:       title,
		Subtitle:    subtitle,
		Tag:         tag,
		ContentMd:   contentMd,
		ContentHtml: contentHtml,
		ContentText: contentText,
		Author:      username.(string),
		AuthorId:    userid.(int64),
		CreateDate:  time.Now(),
	})
	if err != nil {
		c.String(http.StatusOK, "0")
	}
	c.String(http.StatusOK, "1")
}

// ApiArticleEditHandler /api/article PUT
func ApiArticleEditHandler(c *gin.Context) {
	session := sessions.Default(c)
	userid := session.Get("userid")
	username := session.Get("username")
	authority := session.Get("authority")
	if userid == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.String(200, "-1")
		return
	}
	id := c.DefaultPostForm("id", "1")
	title := c.DefaultPostForm("title", "")
	subtitle := c.DefaultPostForm("subtitle", "")
	tag := c.DefaultPostForm("tag", "")
	contentMd := c.DefaultPostForm("content_md", "")
	contentHtml := c.DefaultPostForm("content_html", "")
	reader := strings.NewReader(contentHtml)
	doc, _ := goquery.NewDocumentFromReader(reader)
	contentText := doc.Text()

	err := stores.ArticleStore.UpdateArticle(util.StringToInt64(id), &core.Article{
		Title:        title,
		Subtitle:     subtitle,
		Tag:          tag,
		ContentMd:    contentMd,
		ContentHtml:  contentHtml,
		ContentText:  contentText,
		Author:       username.(string),
		AuthorId:     userid.(int64),
		LastEditDate: time.Now(),
	})
	if err != nil {
		c.String(http.StatusOK, "0")
		return
	}
	c.String(http.StatusOK, "1")
}

// ApiArticleDeleteHandler /api/article DELETE
func ApiArticleDeleteHandler(c *gin.Context) {
	session := sessions.Default(c)
	userid := session.Get("userid")
	authority := session.Get("authority")
	if userid == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.String(200, "-1")
		return
	}
	id := c.DefaultQuery("id", "-1")

	article, _ := stores.ArticleStore.GetArticleById(util.StringToInt64(id))
	trueUserid := article.AuthorId
	if userid.(int64) != trueUserid {
		c.String(200, "-1")
		return
	}

	_ = stores.ArticleStore.DeleteArticle(util.StringToInt64(id))
	_ = stores.CommentStore.DeleteCommentsByArticleId(util.StringToInt64(id))
	c.String(200, "1")
}

// ApiArticleRetrievalHandler /api/article/retrieval POST
func ApiArticleRetrievalHandler(c *gin.Context) {
	ip := c.ClientIP()
	search, _ := stores.SearchStore.GetLatestSearchByIp(ip)
	if search.Ip != "" {
		lastTime := search.SearchDate
		diff := time.Since(lastTime).Seconds()
		if diff < 5.0 {
			c.String(200, "0")
			return
		}
	}

	keyword := c.DefaultPostForm("keyword", "")
	env := util.ParseUserAgent(c.Request.UserAgent())
	go AddSearchRecord(&core.Search{
		Ip:         c.ClientIP(),
		Keyword:    keyword,
		SearchDate: time.Now(),
		Platform:   env.Platform,
		Browser:    env.Browser,
		Version:    env.Version,
	})

	articlesTitle, _ := stores.ArticleStore.GetRelativeArticlesByTitle(keyword, "id", "title", "tag")
	articlesContent, _ := stores.ArticleStore.GetRelativeArticlesByContent(keyword, "id", "title", "tag", "content_text")
	for _, article := range articlesContent {
		if len(article.ContentText) > 300 {
			article.ContentText = article.ContentText[:300]
		}
	}
	c.JSON(200, gin.H{
		"articlesTitle":   articlesTitle,
		"articlesContent": articlesContent,
	})
}

// ApiCommentsAddHandler /api/comments POST
func ApiCommentsAddHandler(c *gin.Context) {
	ip := c.ClientIP()
	comment, _ := stores.CommentStore.GetLatestCommentByIp(ip)

	if comment.Ip != "" {
		lastCommentTime := comment.CommentDate
		diff := time.Since(lastCommentTime).Seconds()
		if diff < 60.0 {
			c.String(200, "2")
			return
		}
	}
	var reviewerId int64
	var replyName string
	session := sessions.Default(c)
	replyName = c.DefaultPostForm("email", "")
	content := c.DefaultPostForm("content", "")
	if content == "" {
		c.String(200, "0")
		return
	}

	userid := session.Get("userid")
	if userid != nil {
		reviewerId = userid.(int64)
		replyName = session.Get("nickname").(string)
	} else {
		if replyName == "" {
			c.String(200, "0")
			return
		}
		reviewerId = -1
	}
	articleId := c.DefaultPostForm("article_id", "-1")
	ua := c.Request.UserAgent()
	env := util.ParseUserAgent(ua)
	location := GetLocation(ip)
	_ = stores.CommentStore.AddComment(&core.Comment{
		ArticleId:   util.StringToInt64(articleId),
		ReplyName:   replyName,
		Content:     content,
		CommentDate: time.Now(),
		Platform:    env.Platform,
		Browser:     env.Browser,
		Version:     env.Version,
		Ip:          ip,
		Location:    location,
		ReviewerId:  reviewerId,
	})
	c.String(200, "1")
}

// ApiCommentsDeleteHandler /api/comments DELETE
func ApiCommentsDeleteHandler(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("username")
	currentAuthority := session.Get("authority")
	if currentUsername == nil || currentAuthority.(int8) != int8(1) {
		c.String(200, "0")
		return
	}

	id := c.DefaultQuery("id", "")
	_ = stores.CommentStore.DeleteCommentById(util.StringToInt64(id))
	c.String(200, "1")
}

func AddSearchRecord(s *core.Search) {
	s.Location = GetLocation(s.Ip)
	_ = stores.SearchStore.AddSearch(s)
}
