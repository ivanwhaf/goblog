package core

import "time"

type Admin struct {
	Id           int64
	Username     string
	Password     string
	Nickname     string
	Sex          string
	Authority    int8
	Avatar       string
	RegisterDate time.Time
}
