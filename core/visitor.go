package core

import "time"

type Visitor struct {
	Id        int64
	Ip        string
	Location  string
	VisitDate time.Time
	Platform  string
	Browser   string
	Version   string
}
