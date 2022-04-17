package accessmodel

type SysPower struct {
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	Title     string     `gorm:"column:title" json:"title"`
	Code      string     `gorm:"column:code" json:"code"`
	Tags      string     `gorm:"column:tags" json:"tags"`
	Desc      string     `gorm:"column:description" json:"desc"`
	Disabled  int        `gorm:"column:disabled" json:"disabled"`
	SysMenuId int        `gorm:"column:sys_menu_id" json:"menuId"`
	SysMenu   *SysMenu   `gorm:"foreignKey:SysMenuId" json:"menu"`
	SysRoles  []*SysRole `gorm:"many2many:sys_role_r_sys_power" json:"roles"`
}

func (power *SysPower) GetResourceName() string {
	return "sysPowers"
}

type SysPowerBody struct {
	Title      *string `json:"title"`
	Code       *string `json:"code"`
	Tags       *string `json:"tags"`
	Desc       *string `json:"desc"`
	Disabled   *int    `json:"disabled"`
	SysMenuId  *int    `json:"menuId"`
	SysRoleIds []int   `json:"roleIds"`
}
