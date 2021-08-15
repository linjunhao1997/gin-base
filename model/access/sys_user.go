package model

import "gin-base/model"

const (
	SysRoles = "SysRoles"
)

type SysUser struct {
	model.Base
	UserName     string     `gorm:"column:username" json:"username"`
	Password     string     `gorm:"column:password" json:"-"`
	SysRoles     []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"sysRoles,omitempty"`
	SysResources `gorm:"-"`
}

// must define
type SysUsers []*SysUser

// must define
func (s SysUsers) GetSize() int {
	return len(s)
}

// must define
func (s SysUsers) GetModel() interface{} {
	return SysUser{}
}
