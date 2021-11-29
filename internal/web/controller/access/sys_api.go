package access

import (
	model "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	"gin-base/internal/web/base"
	"gin-base/internal/web/router"
	"gorm.io/gorm"
)

const (
	SysApiPath = "/sysApis"
)

type SysApiController struct {
	base.Controller
}

func (c *SysApiController) InitController() {
	// create
	router.V1.POST("/sysApis", c.Wrap(func(g *base.Gin) {
		api := &model.SysApi{}
		if ok := g.ValidateJson(api); !ok {
			return
		}

		if err := db.DB.Create(api).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(api, "创建成功")
	}))

	// delete
	router.V1.DELETE("/sysApis/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		api := &model.SysApi{ID: id}
		if err := api.LoadById(); err != nil {
			g.Abort(err)
			return
		}

		db.DB.Transaction(func(tx *gorm.DB) error {

			if err := tx.Delete(&model.SysApi{ID: id}).Error; err != nil {
				return err
			}

			if err := tx.Exec("DELETE FROM sys_role_r_sys_api WHERE sys_api_id = ?", id).Error; err != nil {
				return err
			}

			enforcer := rabc.Enforcer
			if _, err := enforcer.DeletePermission(api.Url, api.Method); err != nil {
				return err
			}

			return enforcer.SavePolicy()
		})

		g.RespSuccess(nil, "删除成功")
	}))

	// update
	router.V1.PATCH("/sysApis/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}

		api := &model.SysApi{ID: id}
		if err := api.LoadById(); err != nil {
			g.Abort(err)
			return
		}

		if ok := g.ValidateJson(api); !ok {
			return
		}
		api.ID = id

		db.DB.Transaction(func(tx *gorm.DB) error {

			if err := db.DB.Save(api).Error; err != nil {
				return err
			}

			if api.Enable < 0 {
				enforcer := rabc.Enforcer
				if _, err := enforcer.DeletePermission(api.Url, api.Method); err != nil {
					return err
				}

				return enforcer.SavePolicy()
			} else if api.Enable > 0 {
				//todo SavePolicy()不知道是否要LoadPolicy, 增加enable字段权限校验
			}

			return nil
		})

		g.RespSuccess(api, "修改成功")
	}))

	// retrieve
	router.V1.GET("/sysApis/:id", c.Wrap(func(g *base.Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		power := &model.SysApi{}
		power.ID = id
		if err := power.LoadById(); err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(power, "")
	}))

	// search
	router.V1.POST("/sysApis/_search", c.Wrap(func(g *base.Gin) {
		param := g.ValidateAllowField(base.NewAllowField("id", "url", "enable"))
		if param == nil {
			return
		}

		apis := make([]*model.SysApi, 0)
		if err := param.Search(db.DB).Find(&apis).Error; err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(param.NewPagination(apis, &model.SysApi{}), "")
	}))

	// all
	router.V1.GET("/sysApis", c.Wrap(func(g *base.Gin) {
		apis := make([]*model.SysApi, 0)
		if err := db.DB.Find(&apis).Error; err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(apis, "")
	}))
}
