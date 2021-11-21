package model

import "gin-base/internal/pkg/db"

type SysPower struct {
	ID        int      `gorm:"column:id;primary_key" json:"id"`
	Title     string   `gorm:"column:title" json:"title"`
	Code      string   `gorm:"column:code" json:"code"`
	Tags      string   `gorm:"column:tags" json:"tags"`
	Desc      string   `gorm:"column:description" json:"desc"`
	Enable    int      `gorm:"column:enable" json:"conditions"`
	SysMenuId int      `gorm:"column:sys_menu_id" json:"menuId"`
	SysMenu   *SysMenu `gorm:"foreignKey:SysMenuId" json:"menu"`
}

func (power *SysPower) LoadById() error {
	err := db.DB.Preload("SysMenu").Find(power, "id = ?", power.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (power *SysPower) Create() error {
	return db.DB.Create(power).Error
}

func (power *SysPower) Update() error {
	return db.DB.Save(power).Error
}

func (power *SysPower) Delete() error {
	return db.DB.Delete(power, "id = ?", power.ID).Error
}
