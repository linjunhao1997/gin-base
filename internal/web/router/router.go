package router

import (
	"gin-base/internal/web/mid"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func G() *gin.Engine {
	return router
}

var V1, V2, SV1, PV1, PV2, AuthV1, AuthV2 *gin.RouterGroup
var jwtMid *jwt.GinJWTMiddleware

type Controller interface {
	InitController()
}

var LoginHandler = func() gin.HandlerFunc {
	return jwtMid.LoginHandler
}

func init() {
	router = gin.New()
	router.Use(gin.Recovery(), mid.LogMiddleware())

	V1 = router.Group("/api/v1")

	jwtMid = mid.NewJwtMiddleware()
	accessHandlerFunc := mid.NewAccessGinHandlerFunc()

	V1.Use(jwtMid.MiddlewareFunc(), accessHandlerFunc)

	V2 = router.Group("/api/v2")
	SV1 = router.Group("/api/s/v1")
	PV1 = router.Group("/api/p/v1")
	PV2 = router.Group("/api/p/v2")
	AuthV1 = router.Group("/auth/v1")
	AuthV2 = router.Group("/auth/v2")
}
