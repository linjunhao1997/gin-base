package model

import (
	"gin-base/internal/pkg/db"
	"gin-base/internal/pkg/rabc"
	gutils "gin-base/internal/utils"
	"gorm.io/gorm"
)

type SysUser struct {
	ID       int        `gorm:"column:id;primary_key" json:"id"`
	Username string     `gorm:"column:username" json:"username"`
	Password string     `gorm:"column:password" json:"password"`
	Phone    string     `gorm:"column:phone" json:"phone"`
	Email    string     `gorm:"column:email" json:"email"`
	Enable   int8       `gorm:"column:enable" json:"conditions"`
	SysRoles []*SysRole `gorm:"many2many:sys_user_r_sys_role" json:"roles"`

	RoleIds []int `gorm:"-" json:"roleIds"`
}

// must define
func (user *SysUser) GetID() int {
	return user.ID
}

// must define
type SysUsers []*SysUser

func (user *SysUser) LoadById() error {
	err := db.DB.Joins(SYSROLES).Find(user, "id = ?", user.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *SysUser) Create() error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		roles := gutils.Int2Strings(SysRoles(user.SysRoles).Ids())

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		enforcer := rabc.Enforcer

		_, err := enforcer.AddRolesForUser(gutils.Int2String(user.ID), roles)
		if err != nil {
			return err
		}

		return nil
	})

}

func (user *SysUser) Update() error {
	return db.DB.Save(user).Error
}

func (user *SysUser) Delete() error {
	return db.DB.Delete(user, "id = ?", user.ID).Error
}

func (user *SysUser) GetApis() (SysApis, error) {
	if err := user.LoadById(); err != nil {
		return SysApis{}, err
	}
	apis := SysRoles(user.SysRoles).DistinctSysApis()
	return apis, nil
}
