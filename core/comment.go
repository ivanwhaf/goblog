package core

import "time"

type Comment struct {
	Id          int64
	ArticleId   int64
	ReplyName   string
	Content     string
	CommentDate time.Time
	Platform    string
	Browser     string
	Version     string
	Ip          string
	Location    string
	ReviewerId  int64
}
