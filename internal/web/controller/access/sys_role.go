package accessapi

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

type SysRoleController struct {
	*base.Controller
}

func (c *SysRoleController) InitController() {
	c.Controller = base.NewController(db.DB, router.V1, &accessmodel.SysRole{})

	c.BuildCreateApi(&accessmodel.SysRoleBody{}, accessservice.CreateRole)

	c.BuildUpdateApi(&accessmodel.SysRoleBody{}, accessservice.UpdateRole)

	c.BuildDeleteApi(accessservice.DeleteRole)

	router.V1.POST("/sysRoles/_search", c.Wrap(c.SearchSysRoles))

	router.V1.GET("/sysRoles/:id", c.Wrap(c.GetSysRole))

	router.V1.GET("/sysRoles", c.Wrap(func(g *base.Gin) {
		roles := make([]*accessmodel.SysRole, 0)
		if err := db.DB.Preload(accessmodel.SYSMENUS).Preload(accessmodel.SYSAPIS).Preload(accessmodel.SYSMENUS + "." + accessmodel.SYSPOWERS).Find(&roles).Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(roles, "")
	}))
}

func (c *SysRoleController) SearchSysRoles(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name", "disabled"+
		""))
	if param == nil {
		return
	}

	roles := make([]*accessmodel.SysRole, 0)
	if err := param.Search(db.DB, accessmodel.SYSMENUS+"."+accessmodel.SYSPOWERS, accessmodel.SYSPOWERS, accessmodel.SYSAPIS).Find(&roles).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(param.NewPagination(roles, &accessmodel.SysRole{}), "")
}

func (c *SysRoleController) GetSysRole(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	var role accessmodel.SysRole
	if err := db.DB.Where("id = ?", id).Take(&role).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(role, "")
}
