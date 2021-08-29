package access

import (
	"gin-base/pkg/base"
	handler "gin-base/pkg/handler/access"
	"gin-base/pkg/router"
)

const (
	SysResourcePath = "/SysResources"
)

type SysResourceController struct {
	base.Controller
}

func (c *SysResourceController) HandlerConfig() {
	router.V1.POST(SysResourcePath, c.Wrap(handler.CreateSysResource))
	router.V1.GET(SysResourcePath, c.Wrap(handler.GetAllSysResource))
}
