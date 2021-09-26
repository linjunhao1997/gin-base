package access

import (
	"gin-base/global"
	model "gin-base/model/access"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
	service "gin-base/service/access"
	"strconv"
	"strings"
)

const (
	SysResourcePath = "/sysResources"
)

type SysResourceController struct {
	base.Controller
}

func (c *SysResourceController) HandlerConfig() {

	router.V1.POST(SysResourcePath+"/_relatedSubResources", c.Wrap(c.RelatedSubResources))

	router.V1.POST(SysResourcePath+"/_clearSubResources", c.Wrap(c.ClearSubResources))

	router.V1.POST(SysResourcePath+"/_search", c.Wrap(c.SearchSysResources))

	router.V1.GET(SysResourcePath+"/:id", c.Wrap(c.GetSysResource))

	router.V1.GET(SysResourcePath, c.Wrap(c.GetAllSysResource))
}

func (c *SysResourceController) SearchSysResources(g *base.Gin) {

	param := g.ValidateAllowField(base.NewAllowField("id", "name"))
	if param == nil {
		return
	}

	slice := make([]string, 4)
	slice[0] = model.SYSSUBRESOURCES
	for i := 1; i < 4; i++ {
		slice[i] = strings.Join([]string{slice[i-1], model.SYSSUBRESOURCES}, ".")
	}

	if param.Eq == nil {
		param.Eq = make(map[string]interface{})
	}

	resources := make([]model.SysResource, 0)
	db := global.DB.Table("sys_resource").Where("sys_resource.type = ?", model.MODULE).Preload(slice[3])
	if err := param.PreSearch(db).Find(&resources).Error; err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(param.NewPagination(resources, &model.SysResource{}, db), "")
}

func (c *SysResourceController) RelatedSubResources(g *base.Gin) {

	body := &model.SysResource{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err := service.RelatedSubResources(body)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(nil, "关联资源成功")
}

func (c *SysResourceController) ClearSubResources(g *base.Gin) {

	body := &model.SysResource{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err := service.ClearSubResources(body)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(nil, "删除关联资源成功")
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

	g.RespSuccess(nil, "关联资源成功")
}

func (c *SysResourceController) GetAllSysResource(g *base.Gin) {
	query, ok := g.C.GetQuery("type")
	var num int
	if ok {
		t, err := strconv.Atoi(query)
		if err != nil {
			g.Abort(err)
			return
		}
		num = t
	}

	list, err := service.GetAllSysResource(num)
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(list, "查询成功")
}

func (c *SysResourceController) GetSysResource(g *base.Gin) {
	id, ok := g.ValidateId()
	if !ok {
		return
	}

	resource := &model.SysResource{}
	err := global.DB.Preload(model.SYSSUBRESOURCES).Where("id = ?", id).Take(resource).Error
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(resource, "")
}
