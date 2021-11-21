package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

const (
	SysMenuPath = "/sysMenus"
)

type SysMenuController struct {
	base.Controller
}

func (c *SysMenuController) InitController() {
	router.V1.GET(SysMenuPath, c.Wrap(c.ListEnableMenus))

	router.V1.GET(SysMenuPath+"/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &model.SysMenu{}
		menu.ID = id

		if err := db.DB.Debug().Preload(model.SYSPOWERS, "enable = 1").Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu, "")
	}))
	router.V1.GET(SysMenuPath+"/:id/sysPowers", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		menu := &model.SysMenu{}
		menu.ID = id

		if err := db.DB.Debug().Preload(model.SYSPOWERS, "enable = 1").Find(&menu).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(menu.SysPowers, "")
	}))
}

func (c *SysMenuController) ListEnableMenus(g *base.Gin) {

	menus := make([]model.SysMenu, 0)
	if err := db.DB.Preload(model.SYSPOWERS, "enable = 1").Find(&menus).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(menus, "")
}
