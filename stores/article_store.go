package stores

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"goblog/core"
	"goblog/util"
)

type ArticleStoreInterface interface {
	GetArticlesOrderById() ([]*core.Article, error)
	GetArticlesOrderByIdWithFields(fields ...string) ([]*core.Article, error)
	GetArticlesOrderByIdOffsetLimit(offset int64, limit int64) ([]*core.Article, error)
	GetMostPopularArticles(limit int64) ([]*core.Article, error)
	GetArticleById(id int64) (*core.Article, error)
	GetLatestArticle() (*core.Article, error)
	GetAllArticlesTags() ([]string, error)
	GetRelativeArticlesByTitle(keyword string, fields ...string) ([]*core.Article, error)
	GetRelativeArticlesByContent(keyword string, fields ...string) ([]*core.Article, error)
	GetArticlesCount() (int64, error)
	AddArticle(article *core.Article) error
	UpdateArticleReadCount(id int64) error
	UpdateArticle(id int64, article *core.Article) error
	DeleteArticle(id int64) error
}

type articleStore struct{}

func NewArticleStore() ArticleStoreInterface {
	return &articleStore{}
}

func (*articleStore) GetArticlesOrderById() ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Order("id desc").Find(&articles)
	return articles, nil
}

func (s *articleStore) GetArticlesOrderByIdWithFields(fields ...string) ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Order("id desc").Select(fields).Find(&articles)
	return articles, nil
}

func (*articleStore) GetArticlesOrderByIdOffsetLimit(offset int64, limit int64) ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Order("id desc").Offset(offset).Limit(limit).Find(&articles)
	if articles != nil {
		return articles, nil
	}
	return articles, errors.New("can not find articles")
}

func (*articleStore) GetMostPopularArticles(limit int64) ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Order("read_count desc").Limit(limit).Find(&articles)
	return articles, nil
}

func (*articleStore) GetArticleById(id int64) (*core.Article, error) {
	db := GetMysqlDB()
	var article core.Article
	db.Table("article").Where("id = ?", id).Find(&article)
	return &article, errors.New("not found article")
}

func (*articleStore) GetLatestArticle() (*core.Article, error) {
	db := GetMysqlDB()
	article := core.Article{Id: -1}
	db.Table("article").Last(&article)
	return &article, nil
}

func (*articleStore) GetAllArticlesTags() ([]string, error) {
	db := GetMysqlDB()
	var tags []string
	db.Table("article").Pluck("tag", &tags)
	tags = util.RemoveDuplicates(tags)
	return tags, nil
}

func (*articleStore) GetRelativeArticlesByTitle(keyword string, fields ...string) ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Select(fields).Where("title LIKE ?", "%"+keyword+"%").Order("id desc").Find(&articles)
	return articles, nil
}

func (*articleStore) GetRelativeArticlesByContent(keyword string, fields ...string) ([]*core.Article, error) {
	db := GetMysqlDB()
	var articles []*core.Article
	db.Table("article").Select(fields).Where("content_md LIKE ?", "%"+keyword+"%").Order("id desc").Find(&articles)
	return articles, nil
}

func (*articleStore) GetArticlesCount() (int64, error) {
	db := GetMysqlDB()
	var count int64
	db.Table("article").Count(&count)
	return count, nil
}

func (*articleStore) AddArticle(article *core.Article) error {
	db := GetMysqlDB()
	db.Table("article").Create(article)
	return nil
}

func (*articleStore) UpdateArticleReadCount(id int64) error {
	db := GetMysqlDB()
	db.Table("article").Where("id = ?", id).Update("read_count", gorm.Expr("read_count + ?", 1))
	return nil
}

func (*articleStore) UpdateArticle(id int64, article *core.Article) error {
	db := GetMysqlDB()
	db.Table("article").Where("id = ?", id).Update(article)
	return nil
}

func (*articleStore) DeleteArticle(id int64) error {
	db := GetMysqlDB()
	db.Table("article").Where("id = ?", id).Delete(&core.Article{})
	return nil
}
