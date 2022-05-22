package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoginHandler /login GET
func LoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
