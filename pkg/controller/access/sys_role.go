package access

import (
	"gin-base/global"
	model "gin-base/model/access"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
	service "gin-base/service/access"
)

const (
	SysRolePath = "/sysRoles"
)

type SysRoleController struct {
	base.Controller
}

func (c *SysRoleController) HandlerConfig() {
	router.V1.POST(SysRolePath, c.Wrap(c.CreateSysRole))

	router.V1.POST(SysRolePath+"/_search", c.Wrap(c.SearchSysRoles))

	router.V1.GET(SysRolePath+"/:id", c.Wrap(c.GetSysRole))

	router.V1.POST(SysRolePath+"/_relatedRoleResources", c.Wrap(c.RelatedRoleResources))
}

func (c *SysRoleController) SearchSysRoles(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name"))
	if param == nil {
		return
	}

	roles := make([]model.SysRole, 0)
	if err := param.Search(model.SYSRESOURCES).Find(&roles).Error; err != nil {
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
	if err := global.DB.Preload(model.SYSRESOURCES).Where("id = ?", id).Take(&role).Error; err != nil {
		g.Abort(err)
		return
	}

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
