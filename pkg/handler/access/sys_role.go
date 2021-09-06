package handler

import (
	model "gin-base/model/access"
	"gin-base/pkg/base"
	service "gin-base/service/access"
)

func SearchSysRoles(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name"))
	if param == nil {
		return
	}

	roles, err := service.SearchSysRoles(param.Search(model.SYSRESOURCES))
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(param.NewPagination(roles, &model.SysRole{}), "")
}

func CreateSysRole(g *base.Gin) {

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

func RelatedRoleResources(g *base.Gin) {
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
