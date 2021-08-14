package handler

import (
	"gin-base/pkg/base"
	service "gin-base/service/access"
	"github.com/gin-gonic/gin"
)

func GetSysUser(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	user, err := service.GetSysUser(id)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(user, "")
}

type ResetPasswordParam struct {
	userId          int
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

func ResetPassword(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	body := ResetPasswordParam{}
	if ok := g.ValidateJson(&body); !ok {
		return
	}
	body.userId = id

	g.C.JSON(200, gin.H{
		"data": "success",
	})

}

func SearchSysUsers(g *base.Gin) {

	param := g.ValidateAllowField(base.NewAllowField("id", "username"))
	if param == nil {
		return
	}

	users, err := service.SearchSysUsers(param)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(users, "")
}
