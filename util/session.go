package util

import "github.com/gin-contrib/sessions"

type Session struct {
	UserId *int64
}

func ParseSession(session sessions.Session) {
}
