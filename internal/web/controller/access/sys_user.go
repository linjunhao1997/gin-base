package accessapi

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	accessservice "gin-base/internal/service/access"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
	"github.com/gin-gonic/gin"
)

type SysUserController struct {
	*base.Controller
}

func (c *SysUserController) InitController() {
	c.Controller = base.NewController(db.DB, router.V1, &accessmodel.SysUser{})

	c.GetRouter().GET("/sysUsers/:id", c.Wrap(c.GetSysUser))

	c.GetRouter().POST("/sysUsers/_search", c.Wrap(c.SearchSysUsers))

	c.BuildUpdateApi(&accessmodel.SysUserBody{}, accessservice.UpdateUser)

	c.BuildCreateApi(&accessmodel.SysUserBody{}, accessservice.CreateUser)

	c.BuildDeleteApi(accessservice.DeleteUser)

	c.GetRouter().GET("/sysUsers/self", c.Wrap(c.GetSelf))
}

func (c *SysUserController) GetSysUser(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	user, err := accessservice.GetSysUser(id)
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

func (c *SysUserController) ResetPassword(g *base.Gin) {

	id, ok := g.ValidateId()
	if !ok {
		return
	}

	body := ResetPasswordParam{}
	if ok := g.ValidateStruct(&body); !ok {
		return
	}
	body.userId = id

	g.C.JSON(200, gin.H{
		"data": "success",
	})

}

func (c *SysUserController) SearchSysUsers(g *base.Gin) {

	param := g.ValidateAllowField(base.NewAllowField("id", "username", "phone", "email", "disabled"))
	if param == nil {
		return
	}

	users := make(accessmodel.SysUsers, 0)
	if err := param.Search(db.DB, accessmodel.SYSROLES).Find(&users).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(param.NewPagination(users, &accessmodel.SysUser{}), "")
}

func (c *SysUserController) GetSelf(g *base.Gin) {
	user := g.EnsureSysUser()

	if user == nil {
		g.RespUnauthorized("")
		return
	}

	if err := db.DB.Debug().Preload("SysRoles.SysMenus", "disabled != 1").Preload("SysRoles.SysPowers", "disabled != 1").Find(&user).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(user, "")
}

func (c *SysUserController) DeleteSysUser(g *base.Gin) {
	id, ok := g.ValidateId()
	if !ok {
		return
	}

	user := &accessmodel.SysUser{}
	user.ID = id
	// 需要移到service里，删除操作涉及多个表
	err := user.Delete()
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(user, "")
}
