package core

import "time"

type Search struct {
	Id         int64
	Ip         string
	Location   string
	Keyword    string
	SearchDate time.Time
	Platform   string
	Browser    string
	Version    string
}
