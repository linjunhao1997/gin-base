package model

import (
	"gin-base/internal/model/common"
	"gin-base/internal/pkg/db"
)

type SysUser struct {
	common.Model
	Username string     `gorm:"column:username" json:"username"`
	Password string     `gorm:"column:password" json:"-"`
	Phone    string     `gorm:"column:phone" json:"phone"`
	Email    string     `gorm:"column:email" json:"email"`
	Enable   int8       `gorm:"column:enable" json:"conditions"`
	SysRoles []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"roles"`
}

// must define
func (user *SysUser) GetID() int {
	return user.ID
}

// must define
type SysUsers []SysUser

func (user *SysUser) LoadById() error {
	err := db.DB.Find(user, "id = ?", user.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *SysUser) Create() error {
	return db.DB.Create(user).Error
}

func (user *SysUser) Update() error {
	return db.DB.Save(user).Error
}

func (user *SysUser) Delete() error {
	return db.DB.Delete(user, "id = ?", user.ID).Error
}
