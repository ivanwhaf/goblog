package handlers

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/core"
	"goblog/services"
	"goblog/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// FilesHandler /files
func FilesHandler(c *gin.Context) {
	var publicFiles []map[string]string
	var privateFiles []map[string]string
	var albumFiles []map[string]string
	var login int
	cfg := config.GetConfig()

	session := sessions.Default(c)
	if session.Get("username") != nil && session.Get("authority") == int8(1) {
		login = 1
		privateFiles = services.GetFiles(cfg.File.PrivateFilePath)
		albumFiles = services.GetFiles(cfg.File.AlbumCompressFilePath)
	}

	publicFiles = services.GetFiles(cfg.File.PublicFilePath)
	publicFilesCount := len(publicFiles)

	c.HTML(http.StatusOK, "files.html", gin.H{
		"publicFiles":          publicFiles,
		"privateFiles":         privateFiles,
		"albumFiles":           albumFiles,
		"login":                login,
		"publicFilesCount":     publicFilesCount,
		"uploadPermission":     cfg.File.UploadPermission,
		"downloadPermission":   cfg.File.DownloadPermission,
		"publicFileAllowType":  cfg.File.PublicFileAllowType,
		"privateFileAllowType": cfg.File.PrivateFileAllowType,
		"albumFileAllowType":   cfg.File.AlbumFileAllowType,
	})
}

// FilesPublicHandler /files/public/:filename GET
func FilesPublicHandler(c *gin.Context) {
	fileName := c.Param("filename")
	session := sessions.Default(c)
	cfg := config.GetConfig()

	if cfg.File.DownloadPermission == false && session.Get("username") == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	path := cfg.File.PublicFilePath + fileName
	if !util.ExistsPath(path) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	if util.IsDir(path) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	env := util.ParseUserAgent(c.Request.UserAgent())
	go services.AddDownloadRecord(&core.Download{
		Filename:     fileName,
		Ip:           c.ClientIP(),
		DownloadDate: time.Now(),
		Platform:     env.Platform,
		Browser:      env.Browser,
		Version:      env.Version,
	})
	c.File(config.GetConfig().File.PublicFilePath + fileName)
}

// FilesPrivateHandler /files/private/:filename GET
func FilesPrivateHandler(c *gin.Context) {
	fileName := c.Param("filename")
	c.File(config.GetConfig().File.PrivateFilePath + fileName)
}

// FilesAlbumHandler /files/album/:filename GET
func FilesAlbumHandler(c *gin.Context) {
	fileName := c.Param("filename")
	c.File(config.GetConfig().File.AlbumCompressFilePath + fileName)
}

// FilesAvatarHandler /files/avatar/:filename/:r
func FilesAvatarHandler(c *gin.Context) {
	fileName := c.Param("filename")
	c.File(config.GetConfig().File.AvatarFilePath + fileName)
}

// ApiFilesPublicUploadHandler /api/files/public POST
func ApiFilesPublicUploadHandler(c *gin.Context) {
	cfg := config.GetConfig()
	if !cfg.File.UploadPermission {
		c.String(http.StatusInternalServerError, "no upload permission")
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusInternalServerError, "get file error")
		return
	}

	name := file.Filename
	s := strings.Split(name, ".")
	ext := s[0]
	if len(s) > 1 {
		ext = s[len(s)-1]
	}

	if !util.ElementInSlice(ext, cfg.File.PublicFileAllowType) {
		c.String(http.StatusInternalServerError, "file ext not allowed")
		return
	}

	err = c.SaveUploadedFile(file, cfg.File.PublicFilePath+name)
	if err != nil {
		fmt.Println("save upload file error", err)
		c.String(http.StatusInternalServerError, "save upload file error")
		return
	}
	fmt.Println("shit")
	env := util.ParseUserAgent(c.Request.UserAgent())
	go services.AddUploadRecord(&core.Upload{
		Filename:   name,
		Ip:         c.ClientIP(),
		UploadDate: time.Now(),
		Platform:   env.Platform,
		Browser:    env.Browser,
		Version:    env.Version,
	})
	c.String(http.StatusOK, "1")
}

// ApiFilesPrivateUploadHandler /api/files/private POST
func ApiFilesPrivateUploadHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.String(http.StatusInternalServerError, "no permission")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("get form file error", err)
		c.String(http.StatusInternalServerError, "get file error")
		return
	}
	name := file.Filename
	s := strings.Split(name, ".")
	ext := s[0]
	if len(s) > 1 {
		ext = s[len(s)-1]
	}

	cfg := config.GetConfig()
	if !util.ElementInSlice(ext, cfg.File.PrivateFileAllowType) {
		c.String(http.StatusInternalServerError, "file ext not allowed")
		return
	}

	err = c.SaveUploadedFile(file, cfg.File.PrivateFilePath+name)
	if err != nil {
		fmt.Println("save upload file error", err)
		c.String(http.StatusInternalServerError, "save upload file error")
		return
	}
	c.String(http.StatusOK, "1")
}

// ApiFilesAlbumUploadHandler /api/files/album POST
func ApiFilesAlbumUploadHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.String(http.StatusInternalServerError, "no permission")
		return
	}

	cfg := config.GetConfig()

	form, err := c.MultipartForm()
	if err != nil {
		c.String(200, "0")
		return
	}
	files := form.File["file"]
	for _, file := range files {
		fileName := file.Filename
		s := strings.Split(fileName, ".")
		ext := s[0]
		if len(s) > 1 {
			ext = s[len(s)-1]
		}
		if !util.ElementInSlice(ext, cfg.File.AlbumFileAllowType) {
			c.String(200, "0")
			return
		}

		dir, err := ioutil.ReadDir(config.GetConfig().File.AlbumRawFilePath)
		if err != nil {
			fmt.Println("read album dir error", err)
			c.String(200, "0")
			return
		}
		ymd := time.Now().String()[:10]
		i := 0
		for _, fi := range dir {
			if !fi.IsDir() {
				if len(fi.Name()) >= 10 && fi.Name()[:10] == ymd {
					i++
				}
			}
		}

		fileName = ymd + "-" + strconv.Itoa(i+1) + "." + ext

		// save raw image
		err = c.SaveUploadedFile(file, cfg.File.AlbumRawFilePath+fileName)
		if err != nil {
			fmt.Println("save upload file error")
			c.String(http.StatusInternalServerError, "0")
			return
		}
		// save compress image
		if util.ElementInSlice(ext, []string{"png", "webp", "gif"}) {
			//f, err := os.Open(cfg.File.AlbumRawFilePath + fileName)
			//if err != nil {
			//	c.String(http.StatusInternalServerError, "0")
			//	return
			//}
			//fmt.Println(f)
			//jpgImg, _ := jpeg.Decode(f)
			//
			//newFile, _ := os.Create(cfg.File.AlbumCompressFilePath + ymd + "-" + strconv.Itoa(i+1) + "." + "jpg")
			//
			//_ = jpeg.Encode(newFile, jpgImg, &jpeg.Options{Quality: 40})
			err = c.SaveUploadedFile(file, cfg.File.AlbumCompressFilePath+fileName)
			if err != nil {
				fmt.Println("save upload file error")
				c.String(http.StatusInternalServerError, "0")
				return
			}
		} else {
			err = c.SaveUploadedFile(file, cfg.File.AlbumCompressFilePath+fileName)
			if err != nil {
				fmt.Println("save upload file error")
				c.String(http.StatusInternalServerError, "0")
				return
			}
		}
	}
	c.String(http.StatusOK, "1")
}

// ApiEditorMdAlbumHandler /api/editormd/album POST
func ApiEditorMdAlbumHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || !util.ElementInSlice(authority, []int8{1, 2}) {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，权限不足",
		})
		return
	}

	file, err := c.FormFile("editormd-image-file")
	if err != nil {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，未发现文件",
		})
		return
	}

	fileName := file.Filename
	s := strings.Split(fileName, ".")
	ext := s[0]
	if len(s) > 1 {
		ext = s[len(s)-1]
	}

	cfg := config.GetConfig()
	if !util.ElementInSlice(ext, cfg.File.AlbumFileAllowType) {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，未经允许的类型",
		})
		return
	}

	dir, err := ioutil.ReadDir(config.GetConfig().File.AlbumRawFilePath)
	if err != nil {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，读取目录失败",
		})
		return
	}
	ymd := time.Now().String()[:10]
	i := 0
	for _, fi := range dir {
		if !fi.IsDir() {
			if len(fi.Name()) >= 10 && fi.Name()[:10] == ymd {
				i++
			}
		}
	}
	fileName = ymd + "-" + strconv.Itoa(i+1) + "." + ext

	err = c.SaveUploadedFile(file, cfg.File.AlbumRawFilePath+fileName)
	if err != nil {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，保存文件失败",
		})
		return
	}
	err = c.SaveUploadedFile(file, cfg.File.AlbumCompressFilePath+fileName)
	if err != nil {
		c.JSON(200, gin.H{
			"success": 0,
			"message": "上传失败，保存文件失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"success": 1,
		"message": "图片上传成功",
		"url":     "/files/album/" + fileName,
	})
}

// ApiFilesUploadPermissionHandler /api/files/upload_permission PUT
func ApiFilesUploadPermissionHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || authority.(int8) != int8(1) {
		c.String(200, "-1")
		return
	}
	perm := c.DefaultQuery("perm", "true")
	if perm == "true" {
		config.GetConfig().File.UploadPermission = true
		c.String(200, "1")
		return
	} else if perm == "false" {
		config.GetConfig().File.UploadPermission = false
		c.String(200, "0")
		return
	}
	c.String(200, "-1")
}

// ApiFilesDownloadPermissionHandler /api/files/download_permission PUT
func ApiFilesDownloadPermissionHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || authority.(int8) != int8(1) {
		fmt.Println(authority)
		c.String(200, "-1")
		return
	}

	perm := c.DefaultQuery("perm", "true")
	if perm == "true" {
		config.GetConfig().File.DownloadPermission = true
		c.String(200, "1")
		return
	} else if perm == "false" {
		config.GetConfig().File.DownloadPermission = false
		c.String(200, "0")
		return
	} else {
		c.String(200, "-1")
		return
	}
}

// ApiFilesPublicDeleteHandler /api/files/public DELETE
func ApiFilesPublicDeleteHandler(c *gin.Context) {
	fileName := c.DefaultQuery("filename", "")
	session := sessions.Default(c)
	cfg := config.GetConfig()

	if session.Get("username") == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	path := cfg.File.PublicFilePath + fileName

	if !util.ExistsPath(path) {
		c.String(200, "0")
		return
	}
	if util.IsDir(path) {
		c.String(200, "0")
		return
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println("file delete error", err)
		c.String(200, "0")
		return
	}
	c.String(200, "1")
}

// ApiFilesPrivateDeleteHandler /api/files/private DELETE
func ApiFilesPrivateDeleteHandler(c *gin.Context) {
	fileName := c.DefaultQuery("filename", "")
	session := sessions.Default(c)
	cfg := config.GetConfig()

	if session.Get("username") == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	path := cfg.File.PrivateFilePath + fileName
	if !util.ExistsPath(path) {
		c.String(200, "0")
		return
	}
	if util.IsDir(path) {
		c.String(200, "0")
		return
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println("file delete error", err)
		c.String(200, "0")
		return
	}
	c.String(200, "1")
}

// ApiFilesAlbumDeleteHandler /api/files/album DELETE
func ApiFilesAlbumDeleteHandler(c *gin.Context) {
	fileName := c.DefaultQuery("filename", "")
	session := sessions.Default(c)
	cfg := config.GetConfig()

	if session.Get("username") == nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	path := cfg.File.AlbumCompressFilePath + fileName
	if !util.ExistsPath(path) {
		c.String(200, "0")
		return
	}
	if util.IsDir(path) {
		c.String(200, "0")
		return
	}
	err := os.Remove(path)
	if err != nil {
		fmt.Println("file delete error", err)
		c.String(200, "0")
		return
	}
	c.String(200, "1")
}
