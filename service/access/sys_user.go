package service

import (
	"gin-base/component/db"
	model "gin-base/model/access"
)

func GetSysUser(id int) (*model.SysUser, error) {
	var user model.SysUser
	if err := db.RABC.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
