package publicapi

import (
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
)

type AuthController struct {
	base.Controller
}

var identityKey = "id"

func (c *AuthController) InitController() {
	router.AuthV1.POST("/login", router.LoginHandler())
}
