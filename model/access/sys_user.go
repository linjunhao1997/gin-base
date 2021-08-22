package model

import (
	"gin-base/model/common"
)

const (
	SysRoles = "SysRoles"
)

type SysUser struct {
	common.Model
	UserName     string     `gorm:"column:username" json:"username"`
	Password     string     `gorm:"column:password" json:"-"`
	SysRoles     []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"sysRoles,omitempty"`
	SysResources `gorm:"-"`
}

// must define
func (user *SysUser) GetID() int {
	return user.ID
}

// must define
type SysUsers []SysUser
