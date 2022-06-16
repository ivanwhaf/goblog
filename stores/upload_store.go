package stores

import "goblog/core"

type UploadStoreInterface interface {
	AddUpload(upload *core.Upload) error
}

type uploadStore struct{}

func NewUploadStore() UploadStoreInterface {
	return &uploadStore{}
}

func (u uploadStore) AddUpload(upload *core.Upload) error {
	db := GetMysqlDB()
	db.Table("upload").Create(upload)
	return nil
}
