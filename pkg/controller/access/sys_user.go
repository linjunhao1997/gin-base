package access

import (
	"gin-base/pkg/base"
	"gin-base/pkg/handler/access"
	"gin-base/pkg/router"
)

const (
	SysUserPath = "/sysUsers"
)

type SysUserController struct {
	base.Controller
}

func (c *SysUserController) HandlerConfig() {

	router.V1.GET(SysUserPath+"/:id", c.Wrap(handler.GetSysUser))

	router.V1.POST(SysUserPath+"/_search", c.Wrap(handler.SearchSysUsers))

	//router.V1.PATCH(SysUserPath+"/:id", handler.ResetPassword)
}
