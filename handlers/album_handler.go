package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/util"
	"io/ioutil"
	"net/http"
)

// AlbumHandler /album GET
func AlbumHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	} else if !util.ElementInSlice(authority, []int8{1, 2}) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	dir, err := ioutil.ReadDir(config.GetConfig().File.AlbumCompressFilePath)
	if err != nil {
		c.HTML(http.StatusOK, "album.html", gin.H{})
		return
	}
	var fileLst []string
	for _, fi := range dir {
		if !fi.IsDir() {
			fileLst = append(fileLst, fi.Name())
		}
	}
	c.HTML(http.StatusOK, "album.html", gin.H{
		"fileLst": fileLst,
	})
}
