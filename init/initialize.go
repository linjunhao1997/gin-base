package init

import (
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	"gin-base/internal/web/controller/access"
	"gin-base/internal/web/controller/public"
	"gin-base/internal/web/mid"
	"gin-base/internal/web/router"
)

func Initialize() {
	db.DB = db.NewDB()
	rabc.Enforcer = rabc.NewCasbinEnforcer(db.DB)
	mid.JwtMiddleware = mid.NewJwtMiddleware(db.DB)

	router.Router = router.NewRouter()
	router.V1.Use(mid.JwtMiddleware.MiddlewareFunc(), mid.CheckAuthByEnforcer(rabc.Enforcer))
	controllers := make([]router.Controller, 0)

	controllers = append(controllers,
		&access.SysUserController{},
		&access.SysRoleController{},
		&access.SysResourceController{},
		&public.AuthController{},
	)

	for _, c := range controllers {
		c.InitController()
	}

}