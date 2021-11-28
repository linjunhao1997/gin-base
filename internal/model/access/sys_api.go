package model

import "gin-base/internal/pkg/db"

type SysApi struct {
	ID     int    `gorm:"column:id;primary_key" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Desc   string `gorm:"column:description" json:"desc"`
	Url    string `gorm:"column:url" json:"url"`
	Method string `gorm:"column:method" json:"method"`
	Enable uint8  `gorm:"column:enable" json:"conditions"`
}

type SysApis []*SysApi

func (api *SysApi) LoadById() error {
	err := db.DB.Find(api, "id = ?", api.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (api *SysApi) Create() error {
	return db.DB.Create(api).Error
}

func (api *SysApi) Update() error {
	return db.DB.Save(api).Error
}

func (api *SysApi) Delete() error {
	return db.DB.Delete(api, "id = ?", api.ID).Error
}
