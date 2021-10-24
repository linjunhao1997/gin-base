package router

import (
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

var V1, V2, PV1, PV2, AuthV1, AuthV2 *gin.RouterGroup

type Controller interface {
	InitController()
}

func NewRouter() *gin.Engine {
	var root = gin.Default()

	V1 = root.Group("/api/v1")

	V2 = root.Group("/api/v2")

	PV1 = root.Group("/api/p/v1")

	PV2 = root.Group("/api/p/v2")

	AuthV1 = root.Group("/auth/v1")

	AuthV2 = root.Group("/auth/v2")

	return root
}
