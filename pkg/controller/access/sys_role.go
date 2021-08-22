package access

import (
	"gin-base/pkg/base"
	handler "gin-base/pkg/handler/access"
	"gin-base/pkg/router"
)

const (
	SysRolePath = "/sysRoles"
)

type SysRoleController struct {
	base.Controller
}

func (c *SysRoleController) HandlerConfig() {
	router.V1.POST(SysRolePath, c.Wrap(handler.CreateSysRole))
	router.V1.POST(SysRolePath+"/_relatedRoleResources", c.Wrap(handler.RelatedRoleResources))
}
