package stores

import (
	"fmt"
	"github.com/go-redis/redis"
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
var RedisStore RedisStoreInterface

var mysqlDB *gorm.DB
var redisDB *redis.Client

func Init() {
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
	RedisStore = NewRedisStore()
	SetupMysqlDB()
	SetupRedisDB()
}

func SetupMysqlDB() {
	cfg := config.GetConfig()
	conn, err := gorm.Open(cfg.MySQL.Name, cfg.MySQL.Args)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to mysql db server.")
	sqlDB := conn.DB()
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.MySQL.ConnMaxLifeTime))
	mysqlDB = conn
}

func GetMysqlDB() *gorm.DB {
	sqlDB := mysqlDB.DB()
	err := sqlDB.Ping()
	if err != nil {
		sqlDB.Close()
		SetupMysqlDB()
	}
	return mysqlDB
}

func SetupRedisDB() {
	cfg := config.GetConfig()
	conn := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	_, err := conn.Ping().Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("connect to redis db server.")
	redisDB = conn

	// load mysql data to redis cache
	fmt.Println("load mysql data to redis.")
	count, _ := VisitorStore.GetVisitorsCount()
	RedisStore.SetVisitorsCount(count)
	fmt.Println("visitors_count:", count)
}

func GetRedisDB() *redis.Client {
	return redisDB
}
