package handler

import (
	"gin-base/handler"
	service "gin-base/service/access"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSysUser(c *gin.Context) {

	id, ok := handler.ValidateId(c)
	if !ok {
		return
	}

	user, err := service.GetSysUser(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"data": user,
	})
}

type ResetPasswordParam struct {
	userId          int
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

func ResetPassword(c *gin.Context) {
	id, ok := handler.ValidateId(c)
	if !ok {
		return
	}

	body := ResetPasswordParam{}
	if ok := handler.ValidateJson(c, body); !ok {
		return
	}
	body.userId = id

	panic("hehe")

	c.JSON(200, gin.H{
		"data": "success",
	})

}
