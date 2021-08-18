package service

import (
	"gin-base/global/db"
	model "gin-base/model/access"
	"gorm.io/gorm"
)

func GetSysUser(id int) (*model.SysUser, error) {
	var user model.SysUser
	if err := db.DB.Preload(model.SysRoles).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func SearchSysUsers(db *gorm.DB) (model.SysUsers, error) {

	sysUsers := make(model.SysUsers, 0)
	if err := db.Find(&sysUsers).Error; err != nil {
		return nil, err
	}
	/*	db.RABC.Debug().Raw(`SELECT * FROM sys_user WHERE username LIKE 'tes%' ORDER BY id desc LIMIT 10`).Scan(&sysUsers)*/
	return sysUsers, nil
}
