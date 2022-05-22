package core

import "time"

type Login struct {
	Id        int64
	Ip        string
	Location  string
	Username  string
	Password  string
	LoginFlag int8
	LoginDate time.Time
	Platform  string
	Browser   string
	Version   string
}
