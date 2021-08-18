package initialize

import (
	"fmt"
	"gin-base/global/db"
	"gin-base/global/rabc"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func CasbinEnforcer() {
	enforcer, err := newCasbin()
	if err != nil {
		panic(fmt.Sprintf("初始化casbins失败: %v", err))
	}
	rabc.Enforcer = enforcer
}

func newCasbin() (*casbin.Enforcer, error) {
	// casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db.DB, "", "sys_casbin")
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("casbin.conf", adapter)
	if err != nil {
		return nil, err
	}

	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return enforcer, nil
}
