package stores

import "goblog/core"

type LoginStoreInterface interface {
	GetLatestLoginByUsername(username string) (*core.Login, error)
	GetLatestSuccessfulLoginByUsername(username string) (*core.Login, error)
	GetLatestLoginByIp(ip string) (*core.Login, error)
	AddLogin(login *core.Login) error
}

type loginStore struct {
}

func NewLoginStore() LoginStoreInterface {
	return &loginStore{}
}

func (l loginStore) GetLatestLoginByUsername(username string) (*core.Login, error) {
	db := GetDB()
	var login core.Login
	db.Table("login").Where("username = ?", username).Last(&login)
	return &login, nil
}

func (l loginStore) GetLatestSuccessfulLoginByUsername(username string) (*core.Login, error) {
	db := GetDB()
	var login core.Login
	db.Table("login").Where("username = ? AND login_flag = 1", username).Last(&login)
	return &login, nil
}

func (l loginStore) GetLatestLoginByIp(ip string) (*core.Login, error) {
	db := GetDB()
	var login core.Login
	db.Table("login").Where("ip = ?", ip).Last(&login)
	return &login, nil
}

func (l loginStore) AddLogin(login *core.Login) error {
	db := GetDB()
	db.Table("login").Create(login)
	return nil
}
