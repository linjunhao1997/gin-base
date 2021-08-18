package model

type SysRole struct {
	ID           int            `gorm:"column:id;primary_key" json:"id"`
	Name         string         `gorm:"column:name" json:"name"`
	Tag          string         `gorm:"column:tag" json:"tag"`
	SysResources []*SysResource `gorm:"many2many:sys_role_r_sys_resource" json:"sysResources,omitempty"`
}
