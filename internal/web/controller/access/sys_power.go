package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/web/base"
	"gin-base/internal/web/param/access"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

type SysPowerController struct {
	base.Controller
}

func (c *SysPowerController) InitController() {
	// create
	router.V1.POST("/sysPowers", c.Wrap(func(g *base.Gin) {
		power := &model.SysPower{}
		if ok := g.ValidateJson(power); !ok {
			return
		}

		power.SysRoles = access.RoleIdsToSysRoles(power.RoleIds)
		if err := db.DB.Create(power).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(power, "创建成功")
	}))

	// delete
	router.V1.DELETE("/sysPowers/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		db.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&model.SysPower{ID: id}).Error; err != nil {
				return err
			}

			if err := tx.Exec("DELETE FROM sys_role_r_sys_power WHERE sys_power_id = ?", id).Error; err != nil {
				return err
			}

			return nil
		})

		g.RespSuccess(nil, "删除成功")
	}))

	// update
	router.V1.PATCH("/sysPowers/:id", c.Wrap(func(g *base.Gin) {

		id, ok := g.ValidateId()
		if !ok {
			return
		}

		power := &model.SysPower{ID: id}
		err := db.DB.Model(power).Take(power).Error
		if err != nil {
			g.Abort(err)
			return
		}

		if ok := g.ValidateJson(power); !ok {
			return
		}
		power.ID = id

		err = db.DB.Transaction(func(tx *gorm.DB) error {
			if power.RoleIds != nil {
				if err := tx.Debug().Model(power).Association(model.SYSROLES).Replace(access.RoleIdsToSysRoles(power.RoleIds)); err != nil {
					return err
				}
			}

			return tx.Save(power).Error

		})

		if err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(power, "更新成功")
	}))

	router.V1.GET("/sysPowers/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		power := &model.SysPower{}
		power.ID = id
		if err := power.LoadById(); err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(power, "")
	}))

	router.V1.POST("/sysPowers/_search", c.Wrap(func(g *base.Gin) {
		param := g.ValidateAllowField(base.NewAllowField("id", "title", "enable"))
		if param == nil {
			return
		}

		powers := make([]*model.SysPower, 0)
		if err := param.Search(db.DB, "SysMenu").Find(&powers).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(param.NewPagination(powers, &model.SysPower{}), "")
	}))
}
