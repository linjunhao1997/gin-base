package init

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	accessservice "gin-base/internal/service/access"
	accessapi "gin-base/internal/web/controller/access"
	publicapi "gin-base/internal/web/controller/public"
	"gin-base/internal/web/mid"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

func Initialize() {
	db.DB = db.NewDB("root:123456@tcp(192.168.100.100:3306)/gorabc?charset=utf8mb4&parseTime=True&loc=Local")

	db.FileDB = db.NewDB("root:123456@tcp(192.168.100.100:3306)/file_meta?charset=utf8mb4&parseTime=True&loc=Local")

	/*err := db.DB.AutoMigrate(&accessmodel.SysUser{}, &accessmodel.SysRole{}, &accessmodel.SysMenu{}, &accessmodel.SysBlock{}, &accessmodel.SysApi{})
	if err != nil {
		panic(err)
	}*/
	rabc.Enforcer = rabc.NewCasbinEnforcer(db.DB)
	mid.JwtMiddleware = mid.NewJwtMiddleware(db.DB)

	router.Router = router.NewRouter()

	router.V1.Use(mid.JwtMiddleware.MiddlewareFunc(), mid.CheckAuthByEnforcer(rabc.Enforcer))
	controllers := make([]router.Controller, 0)

	controllers = append(controllers,
		&accessapi.SysUserController{},
		&accessapi.SysRoleController{},
		&accessapi.SysMenuController{},
		&accessapi.SysPowerController{},
		&accessapi.SysApiController{},
		&publicapi.AuthController{},
	)

	for _, c := range controllers {
		c.InitController()
	}

	routes := router.Router.Routes()
	for _, route := range routes {
		_, err := accessservice.FindApiByUrlAndMethod(route.Path, route.Method)
		if err == gorm.ErrRecordNotFound {
			api := &accessmodel.SysApi{Url: route.Path, Method: route.Method, Disabled: 0}
			if err := db.DB.Save(api).Error; err != nil {
				panic(err)
			}
		}
	}
	//miniov7.MinioCli = miniov7.NewClient(&miniov7.Config{Endpoint: "192.168.100.10:9000", User: "root", Password: "rootroot"})

}
