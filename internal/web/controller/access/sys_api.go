package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

const (
	SysApiPath = "/sysApis"
)

type SysApiController struct {
	base.Controller
}

func (c *SysApiController) InitController() {

	router.V1.GET(SysApiPath+"/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		power := &model.SysApi{}
		power.ID = id
		if err := power.LoadById(); err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(power, "")
	}))

	router.V1.POST(SysApiPath+"/_search", c.Wrap(func(g *base.Gin) {
		param := g.ValidateAllowField(base.NewAllowField("id", "url", "enable"))
		if param == nil {
			return
		}

		apis := make([]*model.SysApi, 0)
		if err := param.Search(db.DB).Find(&apis).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(param.NewPagination(apis, &model.SysApi{}), "")
	}))

	router.V1.GET(SysApiPath, c.Wrap(func(g *base.Gin) {
		apis := make([]*model.SysApi, 0)
		if err := db.DB.Find(&apis, "enable = 1").Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(apis, "")
	}))
}
