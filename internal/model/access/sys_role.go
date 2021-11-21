package model

type SysRole struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Name      string      `gorm:"column:name" json:"title"`
	Desc      string      `gorm:"column:description" json:"desc"`
	Enable    int         `gorm:"column:enable" json:"conditions"`
	SysMenus  []*SysMenu  `gorm:"many2many:sys_role_r_sys_sys_menu" json:"menus"`
	SysPowers []*SysPower `gorm:"many2many:sys_role_r_sys_sys_power" json:"powers"`
	SysApis   []*SysApi   `gorm:"many2many:sys_role_r_sys_sys_Api" json:"apis"`
}
