package stores

import "goblog/core"

type DownloadStoreInterface interface {
	AddDownload(upload *core.Download) error
}

type downloadStore struct{}

func NewDownloadStore() DownloadStoreInterface {
	return &downloadStore{}
}

func (d downloadStore) AddDownload(download *core.Download) error {
	db := GetDB()
	db.Table("download").Create(download)
	return nil
}
