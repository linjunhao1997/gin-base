package model

type SysRole struct {
	ID   int    `gorm:"column:id;primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Tag  string `gorm:"column:tag" json:"tag"`
	//SysUsers `gorm:"many2many:sys_user_r_sys_role"`
	//SysResources
}
