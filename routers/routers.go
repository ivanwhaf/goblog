package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"goblog/handlers"
	"log"
	"net/http"
	"runtime/debug"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("0xffff"))
	r.Use(sessions.Sessions("my", store))

	r.GET("/", handlers.IndexHandler)
	r.POST("/", handlers.WechatHandler)
	r.GET("/about", handlers.AboutHandler)
	r.GET("/archives", handlers.ArchivesHandler)
	r.GET("/tags", handlers.TagsHandler)
	r.GET("/login", handlers.LoginHandler)
	r.GET("/album", handlers.AlbumHandler)
	r.GET("/video", handlers.VideoHandler)
	r.GET("/study", handlers.StudyHandler)
	r.GET("/files", handlers.FilesHandler)
	r.GET("/article/:tag/:id", handlers.ArticleDetailHandler)
	r.GET("/article/add", handlers.ArticleAddHandler)
	r.GET("/article/:tag/:id/edit", handlers.ArticleEditHandler)
	r.GET("/article/search", handlers.ArticleSearchHandler)
	r.GET("/files/public/:filename", handlers.FilesPublicHandler)
	r.GET("/files/private/:filename", handlers.FilesPrivateHandler)
	r.GET("/files/album/:filename", handlers.FilesAlbumHandler)
	r.GET("/files/avatar/:filename", handlers.FilesAvatarHandler)
	r.GET("/admin/:username", handlers.AdminHandler)
	r.GET("/admin/add", handlers.AdminAddHandler)
	r.GET("/admin/:username/edit", handlers.AdminEditHandler)
	r.GET("/manage", handlers.ManageHandler)
	r.GET("/echarts", handlers.EchartsHandler)
	r.GET("/robots.txt", handlers.RobotsHandler)

	// API
	r.POST("/api/article", handlers.ApiArticleAddHandler)
	r.PUT("/api/article", handlers.ApiArticleEditHandler)
	r.DELETE("/api/article", handlers.ApiArticleDeleteHandler)
	r.POST("/api/article/retrieval", handlers.ApiArticleRetrievalHandler)
	r.POST("/api/comments", handlers.ApiCommentsAddHandler)
	r.DELETE("/api/comments", handlers.ApiCommentsDeleteHandler)
	r.POST("/api/admin/authentication", handlers.ApiAdminLoginHandler)
	r.DELETE("/api/admin/authentication", handlers.ApiAdminLogoutHandler)
	r.POST("/api/admin", handlers.ApiAdminAddHandler)
	r.PUT("/api/admin", handlers.ApiAdminEditHandler)
	r.DELETE("/api/admin", handlers.ApiAdminDeleteHandler)
	r.POST("/api/admin/avatar", handlers.ApiAdminAvatarHandler)
	r.POST("/api/files/public", handlers.ApiFilesPublicUploadHandler)
	r.POST("/api/files/private", handlers.ApiFilesPrivateUploadHandler)
	r.POST("/api/files/album", handlers.ApiFilesAlbumUploadHandler)
	r.DELETE("/api/files/public", handlers.ApiFilesPublicDeleteHandler)
	r.DELETE("/api/files/private", handlers.ApiFilesPrivateDeleteHandler)
	r.DELETE("/api/files/album", handlers.ApiFilesAlbumDeleteHandler)
	r.PUT("/api/files/upload_permission", handlers.ApiFilesUploadPermissionHandler)
	r.PUT("/api/files/download_permission", handlers.ApiFilesDownloadPermissionHandler)
	r.POST("/api/editormd/album", handlers.ApiEditorMdAlbumHandler)
	r.POST("/api/study/dictionary/retrieval", handlers.ApiStudyDictionaryRetrievalHandler)
	r.POST("/api/study/dictionary", handlers.ApiStudyDictionaryHandler)
	r.GET("/api/video", handlers.ApiVideoDownloadHandler)
	r.POST("/api/video", handlers.ApiVideoUploadHandler)
	r.PUT("/api/video/status", handlers.ApiVideoStatusHandler)

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	})

	r.Use(func(c *gin.Context) {
		defer func() {
			if re := recover(); re != nil {
				log.Printf("panic: %v\n", re)
				debug.PrintStack()
				c.HTML(200, "500.html", nil)
			}
		}()
		c.Next()
	})
	return r
}
