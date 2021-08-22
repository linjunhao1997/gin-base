package public

import (
	"gin-base/global"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
)

type AuthController struct {
	base.Controller
}

var identityKey = "id"

func (c *AuthController) HandlerConfig() {
	router.AuthV1.POST("/login", global.JwtMiddleware.LoginHandler)
}
