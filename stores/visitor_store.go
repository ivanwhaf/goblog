package stores

import "goblog/core"

type VisitorStoreInterface interface {
	GetVisitorsCount() (int64, error)
	GetRecentVisitors(limit int64) ([]*core.Visitor, error)
	AddVisitor(visitor *core.Visitor) error
}

type visitorStore struct{}

func NewVisitorStore() VisitorStoreInterface {
	return &visitorStore{}
}

func (*visitorStore) GetVisitorsCount() (int64, error) {
	db := GetDB()
	var count int64
	db.Table("visitor").Count(&count)
	return count, nil
}

func (*visitorStore) GetRecentVisitors(limit int64) ([]*core.Visitor, error) {
	db := GetDB()
	var visitors []*core.Visitor
	db.Table("visitor").Order("id desc").Limit(limit).Find(&visitors)
	return visitors, nil
}

func (*visitorStore) AddVisitor(visitor *core.Visitor) error {
	db := GetDB()
	db.Table("visitor").Create(visitor)
	return nil
}
