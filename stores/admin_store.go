package stores

import "goblog/core"

type AdminStoreInterface interface {
	GetAdminById(id int64) (*core.Admin, error)
	GetAdminByUsername(username string) (*core.Admin, error)
	GetAdmins() ([]*core.Admin, error)
	AddAdmin(admin *core.Admin) error
	UpdateAdminById(id int64, admin *core.Admin) error
	UpdateAdminByUsername(username string, admin *core.Admin) error
	DeleteAdminById(id int64) error
}

type adminStore struct {
}

func NewAdminStore() AdminStoreInterface {
	return &adminStore{}
}

func (s *adminStore) GetAdminById(id int64) (*core.Admin, error) {
	db := GetDB()
	var admin core.Admin
	db.Table("admin").Where("id = ?", id).Find(&admin)
	return &admin, nil
}

func (s *adminStore) GetAdminByUsername(username string) (*core.Admin, error) {
	db := GetDB()
	var admin core.Admin
	db.Table("admin").Where("username = ?", username).Find(&admin)
	return &admin, nil
}

func (s *adminStore) GetAdmins() ([]*core.Admin, error) {
	db := GetDB()
	var admins []*core.Admin
	db.Table("admin").Order("id asc").Find(&admins)
	return admins, nil
}

func (s *adminStore) AddAdmin(admin *core.Admin) error {
	db := GetDB()
	db.Table("admin").Create(admin)
	return nil
}

func (s *adminStore) UpdateAdminById(id int64, admin *core.Admin) error {
	db := GetDB()
	db.Table("admin").Where("id = ?", id).Update(admin)
	return nil
}

func (s *adminStore) UpdateAdminByUsername(username string, admin *core.Admin) error {
	db := GetDB()
	db.Table("admin").Where("username = ?", username).Update(admin)
	return nil
}

func (s *adminStore) DeleteAdminById(id int64) error {
	db := GetDB()
	db.Table("admin").Where("id = ?", id).Delete(&core.Admin{})
	return nil
}
