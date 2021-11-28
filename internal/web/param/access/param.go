package access

import model "gin-base/internal/model/access"

type CreateUserParam struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Enable   int8   `json:"conditions" validate:"oneof=-1 1"`
	RoleIds  []int  `json:"roleIds"`
}

type UpdateUserParam = CreateUserParam

func RoleIdsToSysRoles(roleIds []int) []*model.SysRole {
	roles := make([]*model.SysRole, 0)
	for _, id := range roleIds {
		roles = append(roles, &model.SysRole{ID: id})
	}
	return roles
}

type UpdateRoleParam struct {
	Title    string `json:"Title"`
	Enable   int    `json:"conditions" validate:"oneof=-1 1"`
	MenuIds  []int  `json:"menuIds"`
	PowerIds []int  `json:"powerIds"`
	ApiIds   []int  `json:"apiIds"`
}

func MenuIdsToSysMenus(menuIds []int) []*model.SysMenu {
	menus := make([]*model.SysMenu, 0)
	for _, id := range menuIds {
		menus = append(menus, &model.SysMenu{ID: id})
	}
	return menus
}

func PowerIdsToSysPowers(powerIds []int) []*model.SysPower {
	powers := make([]*model.SysPower, 0)
	for _, id := range powerIds {
		powers = append(powers, &model.SysPower{ID: id})
	}
	return powers
}

func ApiIdsToSysApis(apiIds []int) []*model.SysApi {
	apis := make([]*model.SysApi, 0)
	for _, id := range apiIds {
		apis = append(apis, &model.SysApi{ID: id})
	}
	return apis
}
