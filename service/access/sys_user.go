package service

import (
	"gin-base/component/db"
	model "gin-base/model/access"
	"gin-base/pkg/base"
)

func GetSysUser(id int) (*model.SysUser, error) {
	var user model.SysUser
	if err := db.RABC.Preload(model.SysRoles).Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func SearchSysUsers(param *base.SearchParam) ([]*model.SysUser, error) {

	sysUsers := make([]*model.SysUser, 0)
	if err := param.Search(base.NewLoadField(model.SysRoles)).Find(&sysUsers).Error; err != nil {
		return nil, err
	}
	/*	db.RABC.Debug().Raw(`SELECT * FROM sys_user WHERE username LIKE 'tes%' ORDER BY id desc LIMIT 10`).Scan(&sysUsers)*/
	return sysUsers, nil
}
