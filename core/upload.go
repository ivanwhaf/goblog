package core

import "time"

type Upload struct {
	Id         int64
	Filename   string
	Ip         string
	Location   string
	UploadDate time.Time
	Platform   string
	Browser    string
	Version    string
}
