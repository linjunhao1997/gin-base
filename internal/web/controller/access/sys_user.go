package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	service "gin-base/internal/service/access"
	gutils "gin-base/internal/utils"
	"gin-base/internal/web/base"
	"gin-base/internal/web/param/access"
	"gin-base/internal/web/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SysUserController struct {
	base.Controller
}

func (c *SysUserController) InitController() {

	router.V1.GET("/sysUsers/:id", c.Wrap(c.GetSysUser))

	router.V1.POST("/sysUsers/_search", c.Wrap(c.SearchSysUsers))

	router.V1.PATCH("/sysUsers/:id", c.Wrap(c.UpdateSysUser))

	router.V1.GET("/sysUsers/self", c.Wrap(c.GetSelf))

	router.V1.POST("/sysUsers", c.Wrap(func(g *base.Gin) {

		user := &model.SysUser{}

		if ok := g.ValidateJson(user); !ok {
			return
		}

		user.SysRoles = access.RoleIdsToSysRoles(user.RoleIds)

		err := db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Save(user).Error; err != nil {
				return err
			}

			if len(user.RoleIds) > 0 {
				enforcer := rabc.Enforcer
				_, err := enforcer.AddRolesForUser(gutils.Int2String(user.ID), gutils.Int2Strings(user.RoleIds))
				if err != nil {
					return err
				}

				return enforcer.SavePolicy()
			}

			return nil
		})

		if err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(nil, "创建成功")

	}))

	router.V1.DELETE("/sysUsers/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		db.DB.Transaction(func(tx *gorm.DB) error {
			user := &model.SysUser{ID: id}

			if err := tx.Model(user).Association(model.SYSROLES).Clear(); err != nil {
				return err
			}

			if err := tx.Delete(user).Error; err != nil {
				return err
			}

			enforcer := rabc.Enforcer
			if _, err := enforcer.DeleteRolesForUser(gutils.Int2String(user.ID)); err != nil {
				return err
			}

			return enforcer.SavePolicy()
		})

		g.RespSuccess(nil, "删除成功")
	}))
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

	param := g.ValidateAllowField(base.NewAllowField("id", "username", "phone", "email", "enable"))
	if param == nil {
		return
	}

	users := make(model.SysUsers, 0)
	if err := param.Search(db.DB, model.SYSROLES).Find(&users).Error; err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(param.NewPagination(users, &model.SysUser{}), "")
}

func (c *SysUserController) UpdateSysUser(g *base.Gin) {
	id, ok := g.ValidateId()
	if !ok {
		return
	}

	user := &model.SysUser{ID: id}
	err := db.DB.Model(user).Take(user).Error
	if err != nil {
		g.Abort(err)
		return
	}

	if ok := g.ValidateJson(user); !ok {
		return
	}

	err = db.DB.Transaction(func(tx *gorm.DB) error {
		if user.RoleIds != nil {
			if err := tx.Model(user).Association(model.SYSROLES).Replace(access.RoleIdsToSysRoles(user.RoleIds)); err != nil {
				return err
			}

			enforcer := rabc.Enforcer
			// update g
			_, err := enforcer.DeleteRolesForUser(gutils.Int2String(user.ID))
			if err != nil {
				return err
			}
			_, err = enforcer.AddRolesForUser(gutils.Int2String(user.ID), gutils.Int2Strings(user.RoleIds))
			if err != nil {
				return err
			}

			return enforcer.SavePolicy()
		}

		return tx.Omit(model.SYSROLES).Save(user).Error
	})

	if err != nil {
		g.Abort(err)
		return
	}

	g.RespSuccess(user, "更新成功")
}

func (c *SysUserController) GetSelf(g *base.Gin) {
	user := g.EnsureSysUser()

	if user == nil {
		g.RespUnauthorized("")
		return
	}

	user, err := service.GetSysUser(user.ID)
	if err != nil {
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

	user := &model.SysUser{}
	user.ID = id
	// 需要移到service里，删除操作涉及多个表
	err := user.Delete()
	if err != nil {
		g.Abort(err)
		return
	}
	g.RespSuccess(user, "")
}
