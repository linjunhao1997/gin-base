package service

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	gutils "gin-base/internal/utils"
	"strconv"
)

func GetSysUser(id int) (*model.SysUser, error) {
	var user model.SysUser
	if err := db.DB.Preload(model.SYSROLES+"."+model.SYSRESOURCES).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func DisableSysUser(id int) error {
	return db.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("disable", 1).Error
}

func EnableSysUser(id int) error {
	return db.DB.Model(&model.SysUser{}).Where("id = ?", id).Update("disable", 0).Error
}

func RelatedUserRoles(userId int, roleIds ...int) error {
	user := model.SysRole{ID: userId}
	roles := make([]model.SysRole, len(roleIds))
	for i, roleId := range roleIds {
		roles[i] = model.SysRole{ID: roleId}
	}
	return db.DB.Model(&user).Association("SysRoles").Replace(&roles)
}

func UpdatePGRoles(userId int, roleIds ...int) error {

	roles := gutils.Int2String(roleIds)
	user := strconv.Itoa(userId)
	enforcer := rabc.Enforcer
	for _, roleId := range roleIds {
		apis, err := GetApiResources(roleId)
		if err != nil {
			return err
		}
		_, err = enforcer.DeletePermissionsForUser(user)
		if err != nil {
			return err
		}

		rules := make([][]string, 0)
		for _, v := range apis {

			rules = append(rules, []string{strconv.Itoa(roleId), v.Name, v.Tag})
		}
		_, err = enforcer.AddNamedPolicies("p", rules)
		if err != nil {
			return err
		}

	}

	// update g
	_, err := enforcer.DeleteRolesForUser(user)
	if err != nil {
		return err
	}
	_, err = enforcer.AddRolesForUser(user, roles)
	if err != nil {
		return err
	}

	return enforcer.SavePolicy()
}
