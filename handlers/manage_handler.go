package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goblog/config"
	"goblog/services"
	"goblog/stores"
	"goblog/util"
	"net/http"
	"strconv"
)

// ManageHandler /manage GET
func ManageHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	authority := session.Get("authority")
	if username == nil {
		// may cause explorer 301 cache! shitty!
		// better return 302 code!
		c.Redirect(http.StatusFound, "/login")
		return
	} else if !util.ElementInSlice(authority, []int8{1, 2}) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}

	articles, _ := stores.ArticleStore.GetArticlesOrderByIdWithFields("id", "title", "tag", "create_date", "author")
	comments, _ := stores.CommentStore.GetCommentsOrderById()
	var adminsSlc []map[string]string

	if session.Get("authority") == int8(1) {
		admins, _ := stores.AdminStore.GetAdmins()
		for _, admin := range admins {
			m := map[string]string{
				"id":        strconv.FormatInt(admin.Id, 10),
				"username":  admin.Username,
				"password":  admin.Password,
				"nickname":  admin.Nickname,
				"sex":       admin.Sex,
				"authority": strconv.Itoa(int(admin.Authority)),
			}
			latestLoginDate, _ := stores.LoginStore.GetLatestLoginByUsername(admin.Username)
			if latestLoginDate.Ip != "" {
				m["latestLoginDate"] = latestLoginDate.LoginDate.String()
			}
			adminsSlc = append(adminsSlc, m)
		}
	}
	cfg := config.GetConfig()
	publicFiles := services.GetFiles(cfg.File.PublicFilePath)
	privateFiles := services.GetFiles(cfg.File.PrivateFilePath)
	albumFiles := services.GetFiles(cfg.File.AlbumCompressFilePath)

	visitors, _ := stores.VisitorStore.GetRecentVisitors(100)

	c.HTML(http.StatusOK, "manage.html", gin.H{
		"articles":           articles,
		"articlesCount":      len(articles),
		"comments":           comments,
		"commentsCount":      len(comments),
		"admins":             adminsSlc,
		"adminsCount":        len(adminsSlc),
		"publicFiles":        publicFiles,
		"privateFiles":       privateFiles,
		"albumFiles":         albumFiles,
		"visitors":           visitors,
		"uploadPermission":   cfg.File.UploadPermission,
		"downloadPermission": cfg.File.DownloadPermission,
		"authority":          authority,
	})
}
