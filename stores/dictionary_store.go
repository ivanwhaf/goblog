package stores

import "goblog/core"

type DictionaryStoreInterface interface {
	GetRelativeWordsByKey(key string, limit int64) ([]*core.Dictionary, error)
	AddDictionary(dictionary *core.Dictionary) error
}

type dictionaryStore struct{}

func NewDictionaryStore() DictionaryStoreInterface {
	return &dictionaryStore{}
}

func (d dictionaryStore) GetRelativeWordsByKey(key string, limit int64) ([]*core.Dictionary, error) {
	db := GetMysqlDB()
	var dictionary []*core.Dictionary
	db.Table("dictionary").Where("key_ LIKE ?", key+"%").Limit(limit).Find(&dictionary)
	return dictionary, nil
}

func (d dictionaryStore) AddDictionary(dictionary *core.Dictionary) error {
	db := GetMysqlDB()
	db.Table("dictionary").Create(dictionary)
	return nil
}
