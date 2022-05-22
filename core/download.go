package core

import "time"

type Download struct {
	Id           int64
	Filename     string
	Ip           string
	Location     string
	DownloadDate time.Time
	Platform     string
	Browser      string
	Version      string
}
