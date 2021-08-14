package model

const (
	SysRoles = "SysRoles"
)

type SysUser struct {
	ID           int        `gorm:"column:id;primary_key" json:"id"`
	UserName     string     `gorm:"column:username" json:"username"`
	Password     string     `gorm:"column:password" json:"-"`
	SysRoles     []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"sysRoles,omitempty"`
	SysResources `gorm:"-"`
}

type SysUsers []*SysUser
