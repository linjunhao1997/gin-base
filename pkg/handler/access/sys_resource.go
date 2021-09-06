package handler

import (
	model "gin-base/model/access"
	"gin-base/pkg/base"
	service "gin-base/service/access"
	"strings"
)

func SearchSysResources(g *base.Gin) {
	param := g.ValidateAllowField(base.NewAllowField("id", "name"))
	if param == nil {
		return
	}

	slice := make([]string, 5)
	slice[0] = model.SYSRESOURCES
	for i := 1; i < 5; i++ {
		slice[i] = strings.Join([]string{slice[i-1], model.SYSRESOURCES}, ".")
	}

	if param.Eq == nil {
		param.Eq = make(map[string]interface{})
	}
	param.Eq["type"] = model.MODULE
	resources, err := service.SearchSysResources(param.Search(slice...))
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(param.NewPagination(resources, &model.SysResource{}), "")
}

func CreateSysResource(g *base.Gin) {

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

func GetAllSysResource(g *base.Gin) {
	list, err := service.GetAllSysResource()
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(list, "查询成功")
}
