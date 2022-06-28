package accessapi

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

type SysApiController struct {
	*base.Controller
}

func (c *SysApiController) InitController() {
	c.Controller = base.NewController(db.G(), router.V1, &accessmodel.SysApi{})

	c.BuildCreateApi(&accessmodel.SysApiBody{}, accessservice.CreateApi)

	c.BuildRetrieveApi()

	c.BuildDeleteApi(accessservice.DeleteApi)

	c.BuildUpdateApi(&accessmodel.SysApiBody{}, accessservice.UpdateApi)

	c.BuildSearchApi(nil)

	// all
	c.GetRouter().GET("/sysApis", c.Wrap(func(g *base.Gin) {
		apis := make([]*accessmodel.SysApi, 0)
		if err := db.G().Find(&apis).Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(apis, "")
	}))
}
