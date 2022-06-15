package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AboutHandler /about GET
func AboutHandler(c *gin.Context) {
	session := sessions.Default(c)
	c.HTML(http.StatusOK, "about.html", gin.H{
		"username": session.Get("username"),
		"nickname": session.Get("nickname"),
	})
}
