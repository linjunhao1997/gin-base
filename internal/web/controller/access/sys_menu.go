package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

type SysMenuController struct {
	base.Controller
}

func (c *SysMenuController) InitController() {

	// create
	router.V1.POST("/sysMenus", c.Wrap(func(g *base.Gin) {
		menu := &model.SysMenu{}
		if ok := g.ValidateJson(menu); !ok {
			return
		}

		if err := db.DB.Create(menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu, "创建成功")
	}))

	// delete
	router.V1.DELETE("/sysMenus/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&model.SysMenu{ID: id}).Error; err != nil {
				return err
			}

			if err := tx.Exec("DELETE FROM sys_role_r_sys_menu WHERE sys_menu_id = ?", id).Error; err != nil {
				return err
			}

			return nil
		})

		g.RespSuccess(nil, "删除成功")
	}))

	// update
	router.V1.PATCH("/sysMenus/:id", c.Wrap(func(g *base.Gin) {

		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &model.SysMenu{ID: id}
		err := db.DB.Model(menu).Take(menu).Error
		if err != nil {
			g.Abort(err)
			return
		}

		if ok := g.ValidateJson(menu); !ok {
			return
		}
		menu.ID = id

		if err := db.DB.Save(menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu, "更新成功")
	}))

	// retrieve
	router.V1.GET("/sysMenus/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &model.SysMenu{ID: id}
		if err := db.DB.Debug().Preload(model.SYSPOWERS, "enable = 1").Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu, "")
	}))

	router.V1.GET("/sysMenus/:id/sysPowers", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &model.SysMenu{ID: id}
		if err := db.DB.Preload(model.SYSPOWERS + "." + model.SYSROLES).Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu.SysPowers, "")
	}))

	// all
	router.V1.GET("/sysMenus", c.Wrap(c.ListEnableMenus))

}

func (c *SysMenuController) ListEnableMenus(g *base.Gin) {

	menus := make([]model.SysMenu, 0)
	if err := db.DB.Preload(model.SYSPOWERS, "enable = ?", 1).Where("enable = ?", 1).Find(&menus).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(menus, "")
}
