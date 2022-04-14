package accessservice

import (
	accessmodel "gin-base/internal/model/access"
	"gin-base/internal/pkg/db"
)

func GetSysUser(id int) (*accessmodel.SysUser, error) {
	var user accessmodel.SysUser
	user.ID = id

	if err := db.DB.Debug().Preload("SysRoles.SysMenus", "enable = 1").Preload("SysRoles.SysPowers", "enable = 1").Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
