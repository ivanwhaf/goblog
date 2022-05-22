package core

import "time"

type Dictionary struct {
	Id         int64
	Key_       string
	Value_     string
	KeyType    string
	ValueType  string
	CreateDate time.Time
}
