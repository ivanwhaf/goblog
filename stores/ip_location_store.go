package stores

import "goblog/core"

type IpLocationStoreInterface interface {
	GetIpLocationByIp(ip string) (*core.IpLocation, error)
	AddIpLocation(ipLocation *core.IpLocation) error
}

type ipLocationStore struct{}

func NewIpLocationStore() IpLocationStoreInterface {
	return &ipLocationStore{}
}

func (i ipLocationStore) GetIpLocationByIp(ip string) (*core.IpLocation, error) {
	db := GetMysqlDB()
	var ipLocation core.IpLocation
	db.Table("ip_location").Where("ip = ?", ip).Find(&ipLocation)
	return &ipLocation, nil
}

func (i ipLocationStore) AddIpLocation(ipLocation *core.IpLocation) error {
	db := GetMysqlDB()
	db.Table("ip_location").Create(ipLocation)
	return nil
}
