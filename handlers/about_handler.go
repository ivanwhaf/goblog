package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// AboutHandler /about GET
func AboutHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", gin.H{
		"nickname": "",
		"username": "",
	})
}
