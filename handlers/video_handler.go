package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// VideoHandler /video
func VideoHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "video.html", gin.H{})
}
