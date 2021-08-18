package initialize

import (
	"gin-base/pkg/controller/access"
	"gin-base/pkg/controller/public"
	"gin-base/pkg/router"
)

func Load() {
	MySqlGorm()
	CasbinEnforcer()
	JwtMiddleware()
	router.AppendController(&access.SysRoleController{}, &access.SysUserController{}, &public.AuthController{})
	router.ConfigHandler()
}
