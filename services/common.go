package services

import (
	"goblog/core"
	"goblog/stores"
	"goblog/util"
)

func AddLoginRecord(l *core.Login) {
	l.Location = GetLocation(l.Ip)
	_ = stores.LoginStore.AddLogin(l)
}

func AddSearchRecord(s *core.Search) {
	s.Location = GetLocation(s.Ip)
	_ = stores.SearchStore.AddSearch(s)
}

func AddUploadRecord(u *core.Upload) {
	u.Location = GetLocation(u.Ip)
	_ = stores.UploadStore.AddUpload(u)
}

func AddDownloadRecord(d *core.Download) {
	d.Location = GetLocation(d.Ip)
	_ = stores.DownloadStore.AddDownload(d)
}

func AddVisitorRecord(v *core.Visitor) {
	v.Location = GetLocation(v.Ip)
	_ = stores.VisitorStore.AddVisitor(v)
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
