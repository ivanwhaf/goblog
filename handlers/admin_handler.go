package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pochard/commons/randstr"
	"goblog/config"
	"goblog/core"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"strings"
	"time"
)

// AdminHandler /admin/:username GET
func AdminHandler(c *gin.Context) {
	session := sessions.Default(c)
	trueUsername := session.Get("username")
	if trueUsername == nil {
		c.Redirect(302, "/login")
		return
	}
	username := c.Param("username")
	if username != trueUsername.(string) {
		c.HTML(404, "404.html", nil)
		return
	}

	admin, _ := stores.AdminStore.GetAdminByUsername(username)

	c.HTML(http.StatusOK, "admin.html", gin.H{
		"id":        admin.Id,
		"username":  admin.Username,
		"password":  admin.Password,
		"nickname":  admin.Nickname,
		"sex":       admin.Sex,
		"authority": admin.Authority,
		"avatar":    admin.Avatar,
		"random":    randstr.RandomAlphabetic(2),
	})
}

// AdminAddHandler /admin/add GET
func AdminAddHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil {
		c.Redirect(302, "/login")
		return
	}
	if !util.ElementInSlice(authority.(int8), []int8{1}) {
		c.HTML(404, "404.html", nil)
	}
	c.HTML(200, "add_admin.html", nil)
}

// AdminEditHandler /admin/:username/edit GET
func AdminEditHandler(c *gin.Context) {
	session := sessions.Default(c)
	authority := session.Get("authority")
	if session.Get("username") == nil || !util.ElementInSlice(authority.(int8), []int8{1, 2}) {
		c.Redirect(302, "/login")
		return
	}
	username := c.Param("username")
	admin, _ := stores.AdminStore.GetAdminByUsername(username)
	c.HTML(http.StatusOK, "edit_admin.html", gin.H{
		"id":        admin.Id,
		"username":  admin.Username,
		"password":  admin.Password,
		"nickname":  admin.Nickname,
		"sex":       admin.Sex,
		"authority": admin.Authority,
		"avatar":    admin.Avatar,
		"random":    randstr.RandomAlphabetic(2),
	})
}

// AdminAvatarHandler /admin/avatar/:filename/:r GET
func AdminAvatarHandler(c *gin.Context) {
	fileName := c.Param("filename")
	c.File(config.GetConfig().File.AvatarFilePath + fileName)
}

// ApiAdminLoginHandler /api/admin/authentication POST
func ApiAdminLoginHandler(c *gin.Context) {
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	admin, _ := stores.AdminStore.GetAdminByUsername(username)

	if admin.Password != password {
		c.String(http.StatusOK, "0")
		ip := c.ClientIP()
		env := util.ParseUserAgent(c.Request.UserAgent())
		AddLoginRecord(&core.Login{
			Ip:        ip,
			Username:  username,
			Password:  password,
			LoginFlag: 0,
			LoginDate: time.Now(),
			Platform:  env.Platform,
			Browser:   env.Browser,
			Version:   env.Version,
		})
		return
	}

	session := sessions.Default(c)
	if session.Get("username") != nil {
		c.String(http.StatusOK, "-1")
		return
	}
	session.Set("userid", admin.Id)
	session.Set("username", username)
	session.Set("nickname", admin.Nickname)
	session.Set("sex", admin.Sex)
	session.Set("authority", admin.Authority)
	session.Options(sessions.Options{
		Path:   "/",
		MaxAge: config.GetConfig().Admin.SessionMaxAge,
	})
	_ = session.Save()

	ip := c.ClientIP()
	env := util.ParseUserAgent(c.Request.UserAgent())
	AddLoginRecord(&core.Login{
		Ip:        ip,
		Username:  username,
		Password:  password,
		LoginFlag: 1,
		LoginDate: time.Now(),
		Platform:  env.Platform,
		Browser:   env.Browser,
		Version:   env.Version,
	})
	c.String(http.StatusOK, "1")
}

// ApiAdminLogoutHandler /api/admin/authentication DELETE
func ApiAdminLogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("username") == nil {
		c.String(http.StatusOK, "0")
		return
	}
	session.Clear()
	_ = session.Save()
	c.String(http.StatusOK, "1")
}

// ApiAdminAddHandler /api/admin POST
func ApiAdminAddHandler(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("username")
	currentAuthority := session.Get("authority")
	if currentUsername == nil || currentAuthority.(int8) != int8(1) {
		c.String(http.StatusOK, "0")
		return
	}
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	nickname := c.DefaultPostForm("nickname", "")
	sex := c.DefaultPostForm("sex", "")
	authority := c.DefaultPostForm("authority", "")

	if !util.ElementInSlice(sex, []string{"男", "女"}) || len(username) > 30 || len(password) > 30 || len(nickname) > 30 {
		c.String(200, "0")
		return
	}

	_ = stores.AdminStore.AddAdmin(&core.Admin{
		Username:     username,
		Password:     password,
		Nickname:     nickname,
		Sex:          sex,
		Authority:    int8(util.StringToInt64(authority)),
		Avatar:       config.GetConfig().Admin.DefaultAvatarName,
		RegisterDate: time.Now(),
	})

	c.String(http.StatusOK, "1")
}

// ApiAdminEditHandler /api/admin PUT
func ApiAdminEditHandler(c *gin.Context) {
	session := sessions.Default(c)
	usernameS := session.Get("username")
	authorityS := session.Get("authority")
	if usernameS == nil || (authorityS.(int8) != int8(1) && usernameS.(string) != c.Query("username")) {
		c.String(200, "-1")
	}

	id := c.DefaultPostForm("id", "")
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	nickname := c.DefaultPostForm("nickname", "")
	sex := c.DefaultPostForm("sex", "")
	authority := c.DefaultPostForm("authority", "")
	field := c.DefaultPostForm("field", "")
	if field == "nickname" {
		if len(nickname) > 30 {
			c.String(200, "0")
			return
		}
		_ = stores.AdminStore.UpdateAdminById(util.StringToInt64(id), &core.Admin{
			Nickname: nickname,
		})
		c.String(200, "1")
		return
	} else if field == "sex" {
		if len(sex) > 30 {
			c.String(200, "0")
			return
		}
		_ = stores.AdminStore.UpdateAdminById(util.StringToInt64(id), &core.Admin{
			Sex: sex,
		})
		c.String(200, "1")
		return
	} else if field == "password" {
		if len(password) > 30 {
			c.String(200, "0")
			return
		}
		_ = stores.AdminStore.UpdateAdminById(util.StringToInt64(id), &core.Admin{
			Password: password,
		})
		c.String(200, "1")
		return
	} else if field == "all" {
		if authorityS.(int8) != int8(1) {
			c.String(200, "-1")
			return
		}
		if !util.ElementInSlice(sex, []string{"男", "女"}) || len(username) > 30 || len(password) > 30 || len(nickname) > 30 {
			c.String(200, "-1")
			return
		}
		_ = stores.AdminStore.UpdateAdminById(util.StringToInt64(id), &core.Admin{
			Username:  username,
			Password:  password,
			Nickname:  nickname,
			Sex:       sex,
			Authority: int8(util.StringToInt64(authority)),
		})
		c.String(200, "1")
		return
	} else {
		c.String(200, "0")
		return
	}
}

// ApiAdminDeleteHandler /api/admin DELETE
func ApiAdminDeleteHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil || authority.(int8) != int8(1) {
		c.String(http.StatusOK, "0")
	}
	id := c.DefaultQuery("id", "-1")
	_ = stores.AdminStore.DeleteAdminById(util.StringToInt64(id))
	c.String(http.StatusOK, "1")
}

// ApiAdminAvatarHandler /api/admin/avatar POST
func ApiAdminAvatarHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.String(200, "-1")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.String(200, "0")
		return
	}

	fileName := file.Filename
	s := strings.Split(fileName, ".")
	ext := s[0]
	if len(s) > 1 {
		ext = s[len(s)-1]
	}
	if !util.ElementInSlice(ext, config.GetConfig().File.AvatarFileAllowType) {
		c.String(200, "0")
		return
	}
	fileName = username.(string) + "." + ext
	err = c.SaveUploadedFile(file, config.GetConfig().File.AvatarFilePath+fileName)
	if err != nil {
		c.String(200, "0")
		return
	}

	admin, _ := stores.AdminStore.GetAdminByUsername(username.(string))
	admin.Avatar = fileName
	_ = stores.AdminStore.UpdateAdminByUsername(username.(string), admin)

	c.String(http.StatusOK, "1")
}

func AddLoginRecord(l *core.Login) {
	l.Location = GetLocation(l.Ip)
	_ = stores.LoginStore.AddLogin(l)
}
