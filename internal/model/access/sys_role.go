package accessmodel

import (
	"gin-base/internal/pkg/db"
)

type SysRole struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Name      string      `gorm:"column:name" json:"title"`
	Desc      string      `gorm:"column:description" json:"desc"`
	Enable    int         `gorm:"column:enable" json:"conditions"`
	SysMenus  []*SysMenu  `gorm:"many2many:sys_role_r_sys_menu" json:"menus"`
	SysPowers []*SysPower `gorm:"many2many:sys_role_r_sys_power" json:"powers"`
	SysApis   []*SysApi   `gorm:"many2many:sys_role_r_sys_api" json:"apis"`
}

func (role *SysRole) GetResourceName() string {
	return "sysRoles"
}

type SysRoleBody struct {
	Name   *string `json:"title"`
	Desc   *string `json:"desc"`
	Enable *int    `json:"conditions"`

	MenuIds  []int `json:"menuIds"`
	PowerIds []int `json:"powerIds"`
	ApiIds   []int `json:"apiIds"`
}

type SysRoles []*SysRole

func (roles SysRoles) Ids() []int {
	ids := make([]int, 0)
	for _, role := range roles {
		ids = append(ids, role.ID)
	}
	return ids
}

func (role *SysRole) LoadById() error {
	err := db.DB.Model(role).Find(role, "id = ?", role.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (role *SysRole) LoadSysApis() error {
	err := db.DB.Joins(SYSAPIS).Find(role).Error
	if err != nil {
		return err
	}
	return err
}

func (roles SysRoles) DistinctSysApis() SysApis {
	apis := SysApis{}
	mapping := make(map[int]*SysApi)
	for _, role := range roles {
		for _, api := range role.SysApis {
			if mapping[api.ID] == nil {
				mapping[api.ID] = api
				apis = append(apis, api)
			}
		}
	}
	return apis
}
