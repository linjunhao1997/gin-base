package router

import (
	"github.com/gin-gonic/gin"
)

var Root = gin.Default()

var V1 = Root.Group("/api/v1")

var V2 = Root.Group("/api/v2")

var PV1 = Root.Group("/api/p/v1")

var PV2 = Root.Group("/api/p/v2")

var AuthV1 = Root.Group("/auth/v1")

var AuthV2 = Root.Group("/auth/v2")

type Controller interface {
	HandlerConfig()
}

var Controllers = make([]Controller, 0)

func AppendController(other ...Controller) {
	Controllers = append(Controllers, other...)
}

func ConfigHandler() {
	for _, c := range Controllers {
		c.HandlerConfig()
	}
}
