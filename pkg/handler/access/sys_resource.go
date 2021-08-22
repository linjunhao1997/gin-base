package handler

import (
	model "gin-base/model/access"
	"gin-base/pkg/base"
	service "gin-base/service/access"
)

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
