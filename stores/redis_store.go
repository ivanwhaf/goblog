package stores

import (
	"goblog/util"
)

type RedisStoreInterface interface {
	GetVisitorsCount() int64
	SetVisitorsCount(int64)
	IncreaseVisitorsCount()
}

type redisStore struct{}

func NewRedisStore() RedisStoreInterface {
	return &redisStore{}
}

func (r *redisStore) GetVisitorsCount() int64 {
	db := GetRedisDB()
	count := util.StringToInt64(db.Get("visitors_count").Val())
	return count
}

func (r *redisStore) SetVisitorsCount(count int64) {
	db := GetRedisDB()
	db.Set("visitors_count", count, 0)
}

func (r *redisStore) IncreaseVisitorsCount() {
	db := GetRedisDB()
	db.Incr("visitors_count")
}
