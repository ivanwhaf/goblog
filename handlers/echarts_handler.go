package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// EchartsHandler /echarts GET
func EchartsHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "render.html", gin.H{})
}

//func RenderHandler(c *gin.Context) {
//	c.HTML(http.StatusOK, "render.html", gin.H{})
//}
