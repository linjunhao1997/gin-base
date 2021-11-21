package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

const (
	SysPowerPath = "/sysPowers"
)

type SysPowerController struct {
	base.Controller
}

func (c *SysPowerController) InitController() {

	router.V1.GET(SysPowerPath+"/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		power := &model.SysPower{}
		power.ID = id
		if err := power.LoadById(); err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(power, "")
	}))

	router.V1.POST(SysPowerPath+"/_search", c.Wrap(func(g *base.Gin) {
		param := g.ValidateAllowField(base.NewAllowField("id", "title", "enable"))
		if param == nil {
			return
		}

		powers := make([]*model.SysPower, 0)
		if err := param.Search(db.DB, "SysMenu").Find(&powers).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(param.NewPagination(powers, &model.SysPower{}), "")
	}))
}
