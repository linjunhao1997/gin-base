package model

import "gin-base/internal/pkg/db"

type SysMenu struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Title     string      `gorm:"column:title" json:"title"`
	Icon      string      `gorm:"column:icon" json:"icon"`
	Url       string      `gorm:"column:url" json:"url"`
	ParentId  int         `gorm:"column:parent_id" json:"parent"`
	Desc      string      `gorm:"column:description" json:"desc"`
	Enable    int         `gorm:"column:enable" json:"conditions"`
	SysPowers []*SysPower `json:"powers"`
}

func (menu *SysMenu) LoadById() error {
	err := db.DB.Find(menu, "id = ?", menu.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (menu *SysMenu) Create() error {
	return db.DB.Create(menu).Error
}

func (menu *SysMenu) Update() error {
	return db.DB.Save(menu).Error
}

func (menu *SysMenu) Delete() error {
	return db.DB.Delete(menu, "id = ?", menu.ID).Error
}
