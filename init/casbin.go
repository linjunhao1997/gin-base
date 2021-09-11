package initialize

import (
	"fmt"
	"gin-base/global"
	"gin-base/pkg/base"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CasbinEnforcer() {
	enforcer, err := newCasbin()
	if err != nil {
		panic(fmt.Sprintf("初始化casbins失败: %v", err))
	}
	global.Enforcer = enforcer
}

func newCasbin() (*casbin.Enforcer, error) {
	// casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	adapter, err := gormadapter.NewAdapterByDBUseTableName(global.DB, "", "sys_casbin")
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("casbin.conf", adapter)
	if err != nil {
		return nil, err
	}
	enforcer.EnableAutoSave(false)
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		g := base.Gin{C: c}
		sysUser := g.EnsureSysUser()
		if sysUser.UserName == "admin" {
			c.Next()
			return
		}
		ok, err := global.Enforcer.Enforce(strconv.Itoa(sysUser.ID), c.Request.RequestURI, c.Request.Method)
		if err != nil {
			g.Abort(err)
			return
		} else if !ok {
			g.RespForbidden("")
			return

		} else {
			c.Next()
		}
	}
}
