package service

import (
	"gin-base/global"
	model "gin-base/model/access"
)

func CreateSysResource(resource *model.SysResource) error {
	return global.DB.Create(resource).Error
}
