package access

import (
	model "gin-base/model/access"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
	service "gin-base/service/access"
	"strings"
)

const (
	SysResourcePath = "/sysResources"
)

type SysResourceController struct {
	base.Controller
}

func (c *SysResourceController) HandlerConfig() {

	router.V1.POST(SysResourcePath, c.Wrap(c.CreateSysResource))

	router.V1.POST(SysResourcePath+"/_search", c.Wrap(c.SearchSysResources))

	router.V1.GET(SysResourcePath, c.Wrap(c.GetAllSysResource))
}

func (c *SysResourceController) SearchSysResources(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name"))
	if param == nil {
		return
	}

	slice := make([]string, 5)
	slice[0] = model.SYSSUBRESOURCES
	for i := 1; i < 5; i++ {
		slice[i] = strings.Join([]string{slice[i-1], model.SYSSUBRESOURCES}, ".")
	}

	if param.Eq == nil {
		param.Eq = make(map[string]interface{})
	}
	param.Eq["type"] = model.MODULE

	resources := make([]model.SysResource, 0)
	if err := param.Search(slice...).Find(&resources).Error; err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(param.NewPagination(resources, &model.SysResource{}), "")
}

func (c *SysResourceController) CreateSysResource(g *base.Gin) {

	body := &model.SysResource{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err := service.CreateSysResource(body)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(body, "创建资源成功")
}

func (c *SysResourceController) GetAllSysResource(g *base.Gin) {
	list, err := service.GetAllSysResource()
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(list, "查询成功")
}
