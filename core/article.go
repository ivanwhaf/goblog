package core

import "time"

type Article struct {
	Id           int64
	Title        string
	Subtitle     string
	Tag          string
	ContentMd    string
	ContentHtml  string
	ContentText  string
	Author       string
	AuthorId     int64
	ReadCount    int64
	CreateDate   time.Time
	LastEditDate time.Time
}
