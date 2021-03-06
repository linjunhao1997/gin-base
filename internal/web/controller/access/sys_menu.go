package accessapi

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

type SysMenuController struct {
	*base.Controller
}

func (c *SysMenuController) InitController() {

	c.Controller = base.NewController(db.G(), router.V1, &accessmodel.SysMenu{})

	c.BuildCreateApi(&accessmodel.SysMenuBody{}, accessservice.CreateMenu)

	c.BuildDeleteApi(accessservice.DeleteMenu)

	c.BuildUpdateApi(&accessmodel.SysMenuBody{}, accessservice.UpdateMenu)

	// retrieve
	c.GetRouter().GET("/sysMenus/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &accessmodel.SysMenu{ID: id}
		if err := db.G().Debug().Preload(accessmodel.SYSPOWERS, "disabled != 1").Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu, "")
	}))

	c.GetRouter().GET("/sysMenus/:id/sysPowers", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &accessmodel.SysMenu{ID: id}
		if err := db.G().Preload(accessmodel.SYSPOWERS + "." + accessmodel.SYSROLES).Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu.SysPowers, "")
	}))

	// all
	c.GetRouter().GET("/sysMenus", c.Wrap(c.ListAllMenus))

	// sort
	c.GetRouter().POST("/sysMenus/_sort", c.Wrap(func(g *base.Gin) {
		var ids []int
		if err := g.C.ShouldBindJSON(&ids); err != nil {
			g.Abort(err)
			return
		}
		if err := db.G().Transaction(func(tx *gorm.DB) error {
			for i := range ids {
				if err := tx.Model(accessmodel.SysMenu{}).Where("id = ?", ids[i]).Update("sort", i).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return
		}

		g.RespSuccess(nil, "操作成功")
	}))
}

func (c *SysMenuController) ListAllMenus(g *base.Gin) {

	menus := make([]accessmodel.SysMenu, 0)
	if err := db.G().Preload(accessmodel.SYSPOWERS).Find(&menus).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(menus, "")
}
