package db

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var g *DB

type DB struct {
	*gorm.DB
	*casbin.Enforcer
}

func G() *DB {
	return g
}

func (db *DB) GORM() *gorm.DB {
	return g.DB
}

func NewDB(dsn string) (*DB, error) {
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	gormConfig := gorm.Config{
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gormConfig)
	if err != nil {
		return nil, err
	}

	native, err := db.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	native.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	native.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	native.SetConnMaxLifetime(time.Hour)

	enforcer, err := newCasbinEnforcer(db)
	if err != nil {
		return nil, err
	}

	g = &DB{
		db,
		enforcer,
	}

	return g, nil
}

func newCasbinEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
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
