package accessmodel

import (
	"gin-base/internal/pkg/db"
	gutils "gin-base/internal/utils"
	"gorm.io/gorm"
)

type SysUser struct {
	ID        int            `gorm:"column:id;primary_key" json:"id"`
	Username  string         `gorm:"column:username" json:"username"`
	Password  string         `gorm:"column:password" json:"password"`
	Phone     string         `gorm:"column:phone" json:"phone"`
	Email     string         `gorm:"column:email" json:"email"`
	Disabled  int8           `gorm:"column:disabled" json:"disabled"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	SysRoles  []*SysRole     `gorm:"many2many:sys_user_r_sys_role" json:"roles"`
}

func (user *SysUser) GetResourceName() string {
	return "sysUsers"
}

type SysUserBody struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Phone    *string `json:"phone"`
	Email    *string `json:"email"`
	Disabled *int8   `json:"disabled"`
	RoleIds  []int   `json:"roleIds"`
}

// must define
func (user *SysUser) GetID() int {
	return user.ID
}

// must define
type SysUsers []*SysUser

func (user *SysUser) LoadById() error {
	err := db.G().Joins(SYSROLES).Find(user, "id = ?", user.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (user *SysUser) Create() error {
	return db.G().Transaction(func(tx *gorm.DB) error {
		roles := gutils.Int2Strings(SysRoles(user.SysRoles).Ids())

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		_, err := db.G().AddRolesForUser(gutils.Int2String(user.ID), roles)
		if err != nil {
			return err
		}

		return nil
	})

}

func (user *SysUser) Update() error {
	return db.G().Save(user).Error
}

func (user *SysUser) Delete() error {
	return db.G().Delete(user, "id = ?", user.ID).Error
}

func (user *SysUser) GetApis() (SysApis, error) {
	if err := user.LoadById(); err != nil {
		return SysApis{}, err
	}
	apis := SysRoles(user.SysRoles).DistinctSysApis()
	return apis, nil
}
