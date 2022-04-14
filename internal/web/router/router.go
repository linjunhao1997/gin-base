package router

import (
	"gin-base/internal/web/mid"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

var V1, V2, SV1, PV1, PV2, AuthV1, AuthV2 *gin.RouterGroup

type Controller interface {
	InitController()
}

func NewRouter() *gin.Engine {
	var root = gin.New()
	root.Use(gin.Recovery(), mid.LogMiddleware())

	V1 = root.Group("/api/v1")

	V2 = root.Group("/api/v2")

	SV1 = root.Group("/api/s/v1")

	PV1 = root.Group("/api/p/v1")

	PV2 = root.Group("/api/p/v2")

	AuthV1 = root.Group("/auth/v1")

	AuthV2 = root.Group("/auth/v2")

	return root
}
