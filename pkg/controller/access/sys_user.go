package access

import (
	"gin-base/pkg/base"
	"gin-base/pkg/handler/access"
	"gin-base/pkg/router"
)

const (
	SysUserPath = "/SysUsers"
)

type SysUserController struct {
	base.Controller
}

func (c *SysUserController) HandlerConfig() {
	router.V1.GET(SysUserPath+"/:id", c.Wrap(handler.GetSysUser))

	//router.V1.PATCH(SysUserPath+"/:id", handler.ResetPassword)
}

func init() {

}
