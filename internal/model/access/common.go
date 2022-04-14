package accessmodel

const (
	SYSUSERS        = "SysUsers"
	SYSROLES        = "SysRoles"
	SYSMENUS        = "SysMenus"
	SYSPOWERS       = "SysPowers"
	SYSAPIS         = "SysApis"
	SYSRESOURCES    = "SysResources"
	SYSSUBRESOURCES = "SysSubResources"
)

func RoleIdsToSysRoles(roleIds []int) []*SysRole {
	roles := make([]*SysRole, 0)
	for _, id := range roleIds {
		roles = append(roles, &SysRole{ID: id})
	}
	return roles
}

func MenuIdsToSysMenus(menuIds []int) []*SysMenu {
	menus := make([]*SysMenu, 0)
	for _, id := range menuIds {
		menus = append(menus, &SysMenu{ID: id})
	}
	return menus
}

func PowerIdsToSysPowers(powerIds []int) []*SysPower {
	powers := make([]*SysPower, 0)
	for _, id := range powerIds {
		powers = append(powers, &SysPower{ID: id})
	}
	return powers
}

func ApiIdsToSysApis(apiIds []int) []*SysApi {
	apis := make([]*SysApi, 0)
	for _, id := range apiIds {
		apis = append(apis, &SysApi{ID: id})
	}
	return apis
}
