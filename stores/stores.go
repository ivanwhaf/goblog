package stores

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"goblog/config"
	"time"
)

var ArticleStore ArticleStoreInterface
var VisitorStore VisitorStoreInterface
var CommentStore CommentStoreInterface
var UploadStore UploadStoreInterface
var DownloadStore DownloadStoreInterface
var IpLocationStore IpLocationStoreInterface
var AdminStore AdminStoreInterface
var LoginStore LoginStoreInterface
var SearchStore SearchStoreInterface
var DictionaryStore DictionaryStoreInterface

var db *gorm.DB

func Setup() {
	ArticleStore = NewArticleStore()
	VisitorStore = NewVisitorStore()
	CommentStore = NewCommentStore()
	UploadStore = NewUploadStore()
	DownloadStore = NewDownloadStore()
	IpLocationStore = NewIpLocationStore()
	AdminStore = NewAdminStore()
	LoginStore = NewLoginStore()
	SearchStore = NewSearchStore()
	DictionaryStore = NewDictionaryStore()
	SetupDB()
}

func SetupDB() {
	cfg := config.GetConfig()
	conn, err := gorm.Open(cfg.DB.Name, cfg.DB.Args)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to db server")
	sqlDB := conn.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 1)
	db = conn
}

func GetDB() *gorm.DB {
	sqlDB := db.DB()
	err := sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		SetupDB()
	}
	return db
}
