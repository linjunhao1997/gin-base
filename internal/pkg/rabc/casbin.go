package rabc

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func NewCasbinEnforcer(db *gorm.DB) *casbin.Enforcer {
	enforcer, err := NewCasbin(db)
	if err != nil {
		panic(fmt.Sprintf("初始化casbins失败: %v", err))
	}
	return enforcer
}

func NewCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	// casbin不使用事务管理, 因为他内部使用到事务, 重复用会导致冲突
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "", "sys_casbin")
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
