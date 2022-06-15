package stores

import "goblog/core"

type CommentStoreInterface interface {
	GetCommentsByArticleId(articleId int64) ([]*core.Comment, error)
	GetCommentsOrderById() ([]*core.Comment, error)
	GetLatestCommentByArticleId(articleId int64) (*core.Comment, error)
	GetLatestCommentByIp(ip string) (*core.Comment, error)
	AddComment(comment *core.Comment) error
	DeleteCommentById(id int64) error
	DeleteCommentsByArticleId(articleId int64) error
}

type commentStore struct{}

func NewCommentStore() CommentStoreInterface {
	return &commentStore{}
}

func (c commentStore) GetCommentsByArticleId(articleId int64) ([]*core.Comment, error) {
	db := GetDB()
	var comments []*core.Comment
	db.Table("comment").Where("article_id = ?", articleId).Find(&comments)
	return comments, nil
}

func (c commentStore) GetCommentsOrderById() ([]*core.Comment, error) {
	db := GetDB()
	var comments []*core.Comment
	db.Table("comment").Order("id desc").Find(&comments)
	return comments, nil
}

func (c commentStore) GetLatestCommentByArticleId(articleId int64) (*core.Comment, error) {
	db := GetDB()
	var comment core.Comment
	db.Table("comment").Where("article_id = ?", articleId).Last(&comment)
	return &comment, nil
}

func (c commentStore) GetLatestCommentByIp(ip string) (*core.Comment, error) {
	db := GetDB()
	var comment core.Comment
	db.Table("comment").Where("ip = ?", ip).Last(&comment)
	return &comment, nil
}

func (c commentStore) AddComment(comment *core.Comment) error {
	db := GetDB()
	db.Table("comment").Create(&comment)
	return nil
}

func (c commentStore) DeleteCommentById(id int64) error {
	db := GetDB()
	db.Table("comment").Where("id = ?", id).Delete(&core.Comment{})
	return nil
}

func (c commentStore) DeleteCommentsByArticleId(articleId int64) error {
	db := GetDB()
	db.Table("comment").Where("article_id = ?", articleId).Delete(&core.Comment{})
	return nil
}
