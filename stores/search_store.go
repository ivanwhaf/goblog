package stores

import "goblog/core"

type SearchStoreInterface interface {
	GetLatestSearchByIp(ip string) (*core.Search, error)
	AddSearch(search *core.Search) error
}

type searchStore struct{}

func NewSearchStore() SearchStoreInterface {
	return &searchStore{}
}

func (s searchStore) GetLatestSearchByIp(ip string) (*core.Search, error) {
	db := GetDB()
	var search core.Search
	db.Table("search").Where("ip = ?", ip).Last(&search)
	return &search, nil
}

func (s searchStore) AddSearch(search *core.Search) error {
	db := GetDB()
	db.Table("search").Create(search)
	return nil
}
