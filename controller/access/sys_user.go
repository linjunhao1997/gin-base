package access

import (
	"gin-base/handler/access"
	"gin-base/router"
)

const (
	SysUserPath = "/SysUsers"
)

type SysUserController struct {
}

func (c *SysUserController) PathConfig() {
	router.V1.GET(SysUserPath+"/:id", handler.GetSysUser)

	router.V1.PATCH(SysUserPath+"/:id", handler.ResetPassword)
}

func init() {
	router.AppendController(new(SysUserController))
}
