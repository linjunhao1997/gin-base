package accessservice

import (
	accessmodel "gin-base/internal/model/access"
	model "gin-base/internal/model/common"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	gutils "gin-base/internal/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

func CreateApi(data interface{}) (interface{}, error) {
	data = data.(*accessmodel.SysApiBody)
	api := &accessmodel.SysApi{}

	if err := model.Mapping(api, data); err != nil {
		return nil, err
	}

	if err := db.DB.Create(api).Error; err != nil {
		return nil, err
	}

	return api, nil
}

func DeleteApi(id int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		api := &accessmodel.SysApi{ID: id}
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Take(&api).Error; err != nil {
			return err
		}

		if err := tx.Delete(api).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_api WHERE sys_api_id = ?", id).Error; err != nil {
			return err
		}

		enforcer := rabc.Enforcer
		if _, err := enforcer.DeletePermission(api.Url, api.Method); err != nil {
			return err
		}

		return enforcer.SavePolicy()
	})
}

func FindApiByUrlAndMethod(url, method string) (*accessmodel.SysApi, error) {
	api := &accessmodel.SysApi{}
	if err := db.DB.Take(api, "url = ? and method = ?", url, method).Error; err != nil {
		return nil, err
	}
	return api, nil

}

func UpdateApi(id int, data interface{}) (interface{}, error) {
	data = (data).(*accessmodel.SysApiBody)
	api := &accessmodel.SysApi{ID: id}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&api).Error; err != nil {
			return err
		}

		if err := model.Mapping(&api, data); err != nil {
			return err
		}

		if err := tx.Model(api).Debug().Save(&api).Error; err != nil {
			return err
		}

		if api.Enable != 0 {
			enforcer := rabc.Enforcer
			oldPolicies := enforcer.GetFilteredPolicy(1, api.Url, api.Method)
			newPolicies := make([][]string, 0)
			for _, oldPolicy := range oldPolicies {
				newPolicy := make([]string, 0)
				for _, value := range oldPolicy {
					newPolicy = append(newPolicy, value)
				}
				newPolicies = append(newPolicies, newPolicy)
			}

			if api.Enable < 0 {
				for _, policy := range newPolicies {
					policy[len(policy)-1] = "-1"
				}
			} else if api.Enable > 0 {
				for _, policy := range newPolicies {
					policy[len(policy)-1] = "1"
				}
			}
			if _, err := enforcer.UpdatePolicies(oldPolicies, newPolicies); err != nil {
				return err
			}
			return enforcer.SavePolicy()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return api, nil
}

func CreateMenu(data interface{}) (interface{}, error) {
	data = data.(*accessmodel.SysMenuBody)
	menu := &accessmodel.SysMenu{}

	if err := model.Mapping(menu, data); err != nil {
		return menu, nil
	}

	if err := db.DB.Create(menu).Error; err != nil {
		return nil, err
	}

	return menu, nil
}

func UpdateMenu(id int, data interface{}) (interface{}, error) {
	data = data.(*accessmodel.SysMenuBody)
	menu := &accessmodel.SysMenu{ID: id}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&menu).Error; err != nil {
			return err
		}

		if err := model.Mapping(menu, data); err != nil {
			return err
		}

		if err := tx.Model(menu).Debug().Updates(menu).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return menu, nil
}

func DeleteMenu(id int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&accessmodel.SysMenu{ID: id}).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_menu WHERE sys_menu_id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func CreatePower(data interface{}) (interface{}, error) {
	body := data.(*accessmodel.SysPowerBody)
	power := &accessmodel.SysPower{}

	if err := model.Mapping(power, data); err != nil {
		return nil, err
	}

	power.SysRoles = accessmodel.RoleIdsToSysRoles(body.SysRoleIds)
	if err := db.DB.Create(power).Error; err != nil {
		return nil, err
	}

	return power, nil
}

func UpdatePower(id int, data interface{}) (interface{}, error) {
	body := data.(*accessmodel.SysPowerBody)
	power := &accessmodel.SysPower{ID: id}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&power).Error; err != nil {
			return err
		}

		if err := model.Mapping(power, data); err != nil {
			return err
		}

		if body.SysRoleIds != nil {
			if err := tx.Debug().Model(power).Association(accessmodel.SYSROLES).Replace(accessmodel.RoleIdsToSysRoles(body.SysRoleIds)); err != nil {
				return err
			}
		}

		return tx.Save(power).Error

	})
	if err != nil {
		return nil, err
	}

	return power, nil
}

func DeletePower(id int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&accessmodel.SysPower{ID: id}).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_power WHERE sys_power_id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
}

func CreateRole(data interface{}) (interface{}, error) {
	body := data.(*accessmodel.SysRoleBody)
	role := &accessmodel.SysRole{}

	if err := model.Mapping(role, data); err != nil {
		return nil, err
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if body.MenuIds != nil {
			role.SysMenus = accessmodel.MenuIdsToSysMenus(body.MenuIds)
		}
		if body.PowerIds != nil {
			role.SysPowers = accessmodel.PowerIdsToSysPowers(body.PowerIds)
		}
		if body.ApiIds != nil {
			role.SysApis = accessmodel.ApiIdsToSysApis(body.ApiIds)
		}

		if err := tx.Save(role).Error; err != nil {
			return err
		}

		if err := tx.Model(role).Preload(accessmodel.SYSAPIS, "enable = ?", 1).Find(role).Error; err != nil {
			return err
		}

		if len(role.SysApis) > 0 {
			enforcer := rabc.Enforcer

			rules := make([][]string, 0)
			for _, api := range role.SysApis {
				rules = append(rules, []string{strconv.Itoa(role.ID), api.Url, api.Method, strconv.Itoa(api.Enable)})
			}
			_, err := enforcer.AddNamedPolicies("p", rules)
			if err != nil {
				return err
			}

			return enforcer.SavePolicy()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return role, nil
}

func UpdateRole(id int, data interface{}) (interface{}, error) {
	body := (data).(*accessmodel.SysRoleBody)
	role := &accessmodel.SysRole{ID: id}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&role).Error; err != nil {
			return err
		}

		if err := model.Mapping(&role, body); err != nil {
			return err
		}

		if err := tx.Model(role).Debug().Save(&role).Error; err != nil {
			return err
		}

		if body.MenuIds != nil {
			if err := tx.Model(role).Association(accessmodel.SYSMENUS).Replace(accessmodel.MenuIdsToSysMenus(body.MenuIds)); err != nil {
				return err
			}
		}
		if body.PowerIds != nil {
			if err := tx.Model(role).Association(accessmodel.SYSPOWERS).Replace(accessmodel.PowerIdsToSysPowers(body.PowerIds)); err != nil {
				return err
			}
		}
		if body.ApiIds != nil {
			if err := tx.Model(role).Association(accessmodel.SYSAPIS).Replace(accessmodel.ApiIdsToSysApis(body.ApiIds)); err != nil {
				return err
			}
		}

		if err := tx.Preload(accessmodel.SYSMENUS).Preload(accessmodel.SYSPOWERS).Preload(accessmodel.SYSAPIS).Where("id = ?", id).Take(role).Error; err != nil {
			return err
		}

		if err := tx.Omit(accessmodel.SYSMENUS, accessmodel.SYSAPIS, accessmodel.SYSPOWERS).Save(role).Error; err != nil {
			return err
		}

		enforcer := rabc.Enforcer
		_, err := enforcer.DeletePermissionsForUser(strconv.Itoa(role.ID))
		if err != nil {
			return err
		}

		rules := make([][]string, 0)
		for _, api := range role.SysApis {

			rules = append(rules, []string{strconv.Itoa(role.ID), api.Url, api.Method, strconv.Itoa(api.Enable)})
		}
		_, err = enforcer.AddNamedPolicies("p", rules)
		if err != nil {
			return err
		}

		return enforcer.SavePolicy()
	})
	if err != nil {
		return nil, err
	}

	return role, nil
}

func DeleteRole(id int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&accessmodel.SysRole{ID: id}).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_menu WHERE sys_role_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_power WHERE sys_role_id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Exec("DELETE FROM sys_role_r_sys_api WHERE sys_role_id = ?", id).Error; err != nil {
			return err
		}

		enforcer := rabc.Enforcer
		if _, err := enforcer.DeleteRole(gutils.Int2String(id)); err != nil {
			return err
		}

		return enforcer.SavePolicy()
	})
}

func UpdateUser(id int, data interface{}) (interface{}, error) {
	body := (data).(*accessmodel.SysUserBody)
	user := &accessmodel.SysUser{ID: id}
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).Find(&user).Error; err != nil {
			return err
		}
		if err := model.Mapping(user, body); err != nil {
			return err
		}
		if err := tx.Omit(accessmodel.SYSROLES).Save(user).Error; err != nil {
			return err
		}
		if body.RoleIds != nil {
			if err := tx.Model(user).Association(accessmodel.SYSROLES).Replace(accessmodel.RoleIdsToSysRoles(body.RoleIds)); err != nil {
				return err
			}

			enforcer := rabc.Enforcer
			// update g
			_, err := enforcer.DeleteRolesForUser(gutils.Int2String(user.ID))
			if err != nil {
				return err
			}
			_, err = enforcer.AddRolesForUser(gutils.Int2String(user.ID), gutils.Int2Strings(body.RoleIds))
			if err != nil {
				return err
			}

			return enforcer.SavePolicy()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func CreateUser(data interface{}) (interface{}, error) {
	body := data.(*accessmodel.SysUserBody)
	user := &accessmodel.SysUser{}

	if err := model.Mapping(user, data); err != nil {
		return nil, err
	}

	user.SysRoles = accessmodel.RoleIdsToSysRoles(body.RoleIds)

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(user).Error; err != nil {
			return err
		}

		if len(body.RoleIds) > 0 {
			enforcer := rabc.Enforcer
			_, err := enforcer.AddRolesForUser(gutils.Int2String(user.ID), gutils.Int2Strings(body.RoleIds))
			if err != nil {
				return err
			}

			return enforcer.SavePolicy()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(id int) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		user := &accessmodel.SysUser{ID: id}

		if err := tx.Model(user).Association(accessmodel.SYSROLES).Clear(); err != nil {
			return err
		}

		if err := tx.Delete(user).Error; err != nil {
			return err
		}

		enforcer := rabc.Enforcer
		if _, err := enforcer.DeleteRolesForUser(gutils.Int2String(user.ID)); err != nil {
			return err
		}

		return enforcer.SavePolicy()
	})
}
