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
	ID              int           `gorm:"column:id;primary_key" json:"id"`
	Name            string        `gorm:"column:name" json:"name"`
	Tag             string        `gorm:"column:tag" json:"tag"`
	Type            ResourceType  `gorm:"column:type" json:"type"`
	Disable         uint8         `gorm:"disable" json:"disable"`
	SysSubResources []SysResource `gorm:"many2many:sys_resource_r_sys_sub_resource" json:"children,omitempty"`
}

type SysResources []*SysResource
