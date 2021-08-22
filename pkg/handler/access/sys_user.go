package handler

import (
	"gin-base/global"
	model "gin-base/model/access"
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

	users, err := service.SearchSysUsers(param.Search(model.SysRoles))
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(param.NewPagination(users, &model.SysUser{}), "")
}

type UpdateSysRolesParam struct {
	UserId  int   `json:"userId"`
	RoleIds []int `json:"roleIds"`
}

func UpdateSysRoles(g *base.Gin) {
	var body UpdateSysRolesParam
	if ok := g.ValidateJson(&body); !ok {
		return
	}
	var user model.SysUser
	if err := global.DB.Where("id = ?", body.UserId).Take(&user).Error; err != nil {
		g.Abort(err)
		return
	}

	err := service.UpdatePGRoles(user.ID, 1)
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(nil, "更新成功")
}
