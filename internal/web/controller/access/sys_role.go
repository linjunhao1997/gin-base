package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	service "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

const (
	SysRolePath = "/sysRoles"
)

type SysRoleController struct {
	base.Controller
}

func (c *SysRoleController) InitController() {
	router.V1.POST(SysRolePath, c.Wrap(c.CreateSysRole))

	router.V1.POST(SysRolePath+"/_search", c.Wrap(c.SearchSysRoles))

	router.V1.GET(SysRolePath+"/:id", c.Wrap(c.GetSysRole))

	router.V1.POST(SysRolePath+"/_relatedRoleResources", c.Wrap(c.RelatedRoleResources))

	router.V1.GET(SysRolePath, c.Wrap(func(g *base.Gin) {
		roles := make([]*model.SysRole, 0)
		if err := db.DB.Preload(model.SYSMENUS, "enable = 1").Preload(model.SYSAPIS, "enable = 1").Preload(model.SYSMENUS+"."+model.SYSPOWERS, "enable = 1").Find(&roles).Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(roles, "")
	}))

}

func (c *SysRoleController) SearchSysRoles(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name", "enable"))
	if param == nil {
		return
	}

	roles := make([]model.SysRole, 0)
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

	/*sysResource, err := service.GetSysResources(id)
	if err != nil {
		g.Abort(err)
		return
	}*/

	//role.SysResources = sysResource
	g.RespSuccess(role, "")
}

func (c *SysRoleController) CreateSysRole(g *base.Gin) {

	body := &model.SysRole{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err := service.CreateSysRole(body)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(body, "创建角色成功")
}

type RoleResourcesParam struct {
	RoleID      int   `json:"roleId"`
	ResourceIDs []int `json:"resourceIds"`
}

func (c *SysRoleController) RelatedRoleResources(g *base.Gin) {
	body := &RoleResourcesParam{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err := service.RelatedRoleResources(body.RoleID, body.ResourceIDs)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(nil, "角色权限分配成功")
}
