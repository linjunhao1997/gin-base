package service

import (
	"gin-base/global"
	model "gin-base/model/access"
)

func CreateSysResource(resource *model.SysResource) error {
	return global.DB.Create(resource).Error
}

func GetAllSysResource() (model.SysResources, error) {
	list := model.SysResources{}
	err := global.DB.Model(&model.SysResource{}).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil

}
