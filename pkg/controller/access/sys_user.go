package access

import (
	"gin-base/global"
	model "gin-base/model/access"
	"gin-base/pkg/base"
	"gin-base/pkg/router"
	service "gin-base/service/access"
	"github.com/gin-gonic/gin"
)

const (
	SysUserPath = "/sysUsers"
)

type SysUserController struct {
	base.Controller
}

func (c *SysUserController) HandlerConfig() {

	router.V1.GET(SysUserPath+"/:id", c.Wrap(c.GetSysUser))

	router.V1.POST(SysUserPath+"/_search", c.Wrap(c.SearchSysUsers))

	router.V1.POST(SysUserPath+"/_updateRoles", c.Wrap(c.UpdateSysRoles))

	router.V1.PATCH(SysUserPath+"/:id", c.Wrap(c.UpdateSysUser))

	//router.V1.PATCH(SysUserPath+"/:id", handler.ResetPassword)
}

func (c *SysUserController) GetSysUser(g *base.Gin) {

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

func (c *SysUserController) EnableSysUser(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	err := service.EnableSysUser(id)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(nil, "启用成功")
}

func (c *SysUserController) DisableSysUser(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	err := service.DisableSysUser(id)
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(nil, "禁用成功")
}

type ResetPasswordParam struct {
	userId          int
	OldPassword     string `json:"oldPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

func (c *SysUserController) ResetPassword(g *base.Gin) {

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

func (c *SysUserController) SearchSysUsers(g *base.Gin) {

	param := g.ValidateAllowField(base.NewAllowField("id", "username"))
	if param == nil {
		return
	}

	users := make(model.SysUsers, 0)
	if err := param.Search(model.SYSROLES).Find(&users).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(param.NewPagination(users, &model.SysUser{}), "")
}

type UpdateSysUserParam struct {
	Disable int8 `json:"disable" validate:"oneof=0 1"`
}

func (c *SysUserController) UpdateSysUser(g *base.Gin) {
	id, ok := g.ValidateId()
	if !ok {
		return
	}
	user := &model.SysUser{}
	err := global.DB.Where("id = ?", id).Take(user).Error
	if err != nil {
		g.Abort(err)
		return
	}

	body := &UpdateSysUserParam{}
	if ok := g.ValidateJson(body); !ok {
		return
	}

	err = global.DB.Model(user).Select("Disable", "Disable").Updates(model.SysUser{Disable: body.Disable}).Error
	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(user, "更新成功")
}

type UpdateSysRolesParam struct {
	UserId  int   `json:"userId"`
	RoleIds []int `json:"roleIds"`
}

func (c *SysUserController) UpdateSysRoles(g *base.Gin) {
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
