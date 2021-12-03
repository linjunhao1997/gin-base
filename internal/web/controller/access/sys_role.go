package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	gutils "gin-base/internal/utils"
	"gin-base/internal/web/base"
	"gin-base/internal/web/param/access"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
	"strconv"
)

type SysRoleController struct {
	base.Controller
}

func (c *SysRoleController) InitController() {
	router.V1.POST("/sysRoles", c.Wrap(func(g *base.Gin) {
		role := &model.SysRole{}
		if ok := g.ValidateStruct(role); !ok {
			return
		}

		err := db.DB.Transaction(func(tx *gorm.DB) error {
			if role.MenuIds != nil {
				role.SysMenus = access.MenuIdsToSysMenus(role.MenuIds)
			}
			if role.PowerIds != nil {
				role.SysPowers = access.PowerIdsToSysPowers(role.PowerIds)
			}
			if role.ApiIds != nil {
				role.SysApis = access.ApiIdsToSysApis(role.ApiIds)
			}

			if err := tx.Save(role).Error; err != nil {
				return err
			}

			if err := tx.Model(role).Preload(model.SYSAPIS, "enable = ?", 1).Find(role).Error; err != nil {
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
			g.Abort(err)
			return
		}

		g.RespSuccess(role, "更新成功")
	}))

	router.V1.POST("/sysRoles/_search", c.Wrap(c.SearchSysRoles))

	router.V1.GET("/sysRoles/:id", c.Wrap(c.GetSysRole))

	router.V1.GET("/sysRoles", c.Wrap(func(g *base.Gin) {
		roles := make([]*model.SysRole, 0)
		if err := db.DB.Preload(model.SYSMENUS, "enable = 1").Preload(model.SYSAPIS, "enable = 1").Preload(model.SYSMENUS+"."+model.SYSPOWERS, "enable = 1").Find(&roles).Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(roles, "")
	}))

	router.V1.PATCH("/sysRoles/:id", c.Wrap(func(g *base.Gin) {

		id, ok := g.ValidateId()
		if !ok {
			return
		}

		role := &model.SysRole{}
		err := db.DB.Preload(model.SYSMENUS).Preload(model.SYSPOWERS).Preload(model.SYSAPIS).Where("id = ?", id).Take(role).Error
		if err != nil {
			g.Abort(err)
			return
		}

		if ok := g.ValidateStruct(role); !ok {
			return
		}
		role.ID = id
		err = db.DB.Transaction(func(tx *gorm.DB) error {
			if role.MenuIds != nil {
				if err := tx.Model(role).Association(model.SYSMENUS).Replace(access.MenuIdsToSysMenus(role.MenuIds)); err != nil {
					return err
				}
			}
			if role.PowerIds != nil {
				if err := tx.Model(role).Association(model.SYSPOWERS).Replace(access.PowerIdsToSysPowers(role.PowerIds)); err != nil {
					return err
				}
			}
			if role.ApiIds != nil {
				if err := tx.Model(role).Association(model.SYSAPIS).Replace(access.ApiIdsToSysApis(role.ApiIds)); err != nil {
					return err
				}
			}

			if err := tx.Preload(model.SYSMENUS).Preload(model.SYSPOWERS).Preload(model.SYSAPIS).Where("id = ?", id).Take(role).Error; err != nil {
				return err
			}

			if err = tx.Omit(model.SYSMENUS, model.SYSAPIS, model.SYSPOWERS).Save(role).Error; err != nil {
				return err
			}

			enforcer := rabc.Enforcer
			_, err = enforcer.DeletePermissionsForUser(strconv.Itoa(role.ID))
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
			g.Abort(err)
			return
		}

		g.RespSuccess(role, "更新成功")
	}))

	router.V1.DELETE("/sysRoles/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&model.SysRole{ID: id}).Error; err != nil {
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

		g.RespSuccess(nil, "删除成功")
	}))
}

func (c *SysRoleController) SearchSysRoles(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name", "enable"))
	if param == nil {
		return
	}

	roles := make([]*model.SysRole, 0)
	if err := param.Search(db.DB, model.SYSMENUS+"."+model.SYSPOWERS, model.SYSPOWERS, model.SYSAPIS).Find(&roles).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(param.NewPagination(roles, &model.SysRole{}), "")
}

func (c *SysRoleController) GetSysRole(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	var role model.SysRole
	if err := db.DB.Where("id = ?", id).Take(&role).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(role, "")
}
