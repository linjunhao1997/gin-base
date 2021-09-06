package model

type ResourceType uint8

const (
	_         ResourceType = iota
	MODULE                 // 模块
	PAGE                   // 页面
	Component              // 组件
	Button                 // 按钮
	API                    // API
)

type SysResource struct {
	ID           int           `gorm:"column:id;primary_key" json:"id"`
	Name         string        `gorm:"column:name" json:"name"`
	Tag          string        `gorm:"column:tag" json:"tag"`
	Type         ResourceType  `gorm:"column:type" json:"type"`
	Disable      uint8         `gorm:"disable" json:"disable"`
	ParentID     uint8         `gorm:"parent_id" json:"parentId"`
	SysResources []SysResource `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	//SysRoles []*SysRole   `gorm:"many2many:sys_role_r_sys_resource"`
}

type SysResources []*SysResource
