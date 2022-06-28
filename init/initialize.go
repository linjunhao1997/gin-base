package init

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	accessapi "gin-base/internal/web/controller/access"
	publicapi "gin-base/internal/web/controller/public"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

func Initialize() {
	_, err := db.NewDB("root:123456@tcp(192.168.100.100:3306)/gorabc?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

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

	routes := router.G().Routes()
	for _, route := range routes {
		_, err := accessservice.FindApiByUrlAndMethod(route.Path, route.Method)
		if err == gorm.ErrRecordNotFound {
			api := &accessmodel.SysApi{Url: route.Path, Method: route.Method, Disabled: 0}
			if err := db.G().Save(api).Error; err != nil {
				panic(err)
			}
		}
	}
	router.G().Run()
	//miniov7.MinioCli = miniov7.NewClient(&miniov7.Config{Endpoint: "192.168.100.10:9000", User: "root", Password: "rootroot"})

}
