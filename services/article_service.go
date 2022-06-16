package services

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pochard/commons/randstr"
	"goblog/core"
	"goblog/stores"
	"strconv"
	"strings"
)

type Node struct {
	Children []*Node
	Text     string
	Id       string
}

func GenerateContentsNodeLst(html string) []*Node {
	reader := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(reader)

	var contentsNodeLst []*Node
	var lastH1Node *Node
	doc.Find("body").Children().Each(func(i int, selection *goquery.Selection) {
		if selection.Nodes[0].Data == "h1" {
			node := &Node{}
			node.Id, _ = selection.Attr("id")
			node.Text = selection.Text()
			contentsNodeLst = append(contentsNodeLst, node)
			lastH1Node = node
		} else if selection.Nodes[0].Data == "h2" {
			node := &Node{}
			node.Id, _ = selection.Attr("id")
			node.Text = selection.Text()
			if lastH1Node != nil {
				lastH1Node.Children = append(lastH1Node.Children, node)
			} else {
				contentsNodeLst = append(contentsNodeLst, node)
			}
		}
	})
	return contentsNodeLst
}

func GetPrevNextArticleInfo(articles []*core.Article, currentArticleId int64) (string, string, string, string, string, string) {
	var previousId, nextId string
	var previousTag, nextTag string
	var previousTagTitle, nextTagTitle string
	for idx, article := range articles {
		if article.Id == currentArticleId {
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
	return previousId, nextId, previousTag, nextTag, previousTagTitle, nextTagTitle
}

func GetCommentsSlcForArticle(comments []*core.Comment) []map[string]string {
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
	return commentsSlc
}

func GetContentText(html string) string {
	reader := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(reader)
	contentText := doc.Text()
	return contentText
}
