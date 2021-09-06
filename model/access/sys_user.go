package model

import (
	"gin-base/model/common"
)

type SysUser struct {
	common.Model
	UserName     string     `gorm:"column:username" json:"username"`
	Password     string     `gorm:"column:password" json:"-"`
	Disable      int8       `gorm:"column:disable" json:"disable"`
	SysRoles     []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"sysRoles"`
	SysResources `gorm:"-"`
}

// must define
func (user *SysUser) GetID() int {
	return user.ID
}

// must define
type SysUsers []SysUser
