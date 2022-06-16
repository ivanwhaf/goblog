package handlers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var img []byte
var fps = 1.0
var width = 400
var height = 300
var openFlag = true

// VideoHandler /video GET
func VideoHandler(c *gin.Context) {
	session := sessions.Default(c)
	token := c.DefaultQuery("token", "")
	if session.Get("username") == nil && token != "666" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "video.html", gin.H{})
}

// ApiVideoDownloadHandler /api/video GET
func ApiVideoDownloadHandler(c *gin.Context) {
	//c.Header("Cache-Control", "no-cache")
	//c.Header("Pragma", "no-cache")
	//imgs := base64.StdEncoding.EncodeToString(image)
	//c.Data(200, "image/jpeg", []byte(imgs))

	//c.Stream(func(w io.Writer) bool {
	//	c.SSEvent("image", img)
	//	//c.Data(200, "video/webm", img)
	//	time.Sleep(time.Second * 1)
	//	return true
	//})
	token := c.DefaultQuery("token", "")
	session := sessions.Default(c)
	if session.Get("username") == nil && token != "666" {
		c.String(404, "not found")
		return
	}
	if len(img) == 0 {
		c.String(500, "no img")
		return
	}
	w := c.DefaultQuery("width", "400")
	h := c.DefaultQuery("height", "300")
	f := c.DefaultQuery("fps", "1")
	width, _ = strconv.Atoi(w)
	height, _ = strconv.Atoi(h)
	fps, _ = strconv.ParseFloat(f, 32)
	c.Data(200, "image/jpeg", img)
}

// ApiVideoUploadHandler /api/video POST
func ApiVideoUploadHandler(c *gin.Context) {
	file, _ := c.FormFile("file")
	open, err := file.Open()
	if err != nil {
		fmt.Println(err)
		c.String(500, "file open error")
		return
	}

	buf := make([]byte, 10000000)
	n, err := open.Read(buf)
	if err != nil {
		fmt.Println(err)
		c.String(500, "read buf error")
		return
	}

	img = make([]byte, n)
	img = buf[:n]
	c.JSON(200, gin.H{
		"width":     width,
		"height":    height,
		"fps":       fps,
		"open_flag": openFlag,
	})
}

// ApiVideoStatusHandler /api/video/status PUT
func ApiVideoStatusHandler(c *gin.Context) {
	option := c.DefaultQuery("option", "open")
	if option == "open" {
		openFlag = true
	} else {
		openFlag = false
	}
	c.String(200, "1")
}
