package accessapi

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

type SysPowerController struct {
	*base.Controller
}

func (c *SysPowerController) InitController() {

	c.Controller = base.NewController(db.DB, router.V1, &accessmodel.SysPower{})

	c.BuildCreateApi(&accessmodel.SysPowerBody{}, accessservice.CreatePower)

	c.BuildUpdateApi(&accessmodel.SysPowerBody{}, accessservice.UpdatePower)

	c.BuildDeleteApi(accessservice.DeletePower)

	router.V1.GET("/sysPowers/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		power := &accessmodel.SysPower{}
		power.ID = id
		err := db.DB.Preload("SysMenu").Find(power, "id = ?", power.ID).Error
		if err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(power, "查询成功")
	}))

	c.GetRouter().POST("/sysPowers/_search", c.Wrap(func(g *base.Gin) {
		param := g.ValidateAllowField(base.NewAllowField("id", "title", "disabled"))
		if param == nil {
			return
		}

		powers := make([]*accessmodel.SysPower, 0)
		if err := param.Search(db.DB, "SysMenu").Find(&powers).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(param.NewPagination(powers, &accessmodel.SysPower{}), "")
	}))
}
