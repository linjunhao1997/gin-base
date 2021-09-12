package service

import (
	"gin-base/global"
	model "gin-base/model/access"
)

func RelatedSubResources(resource *model.SysResource) error {

	oldResources := make([]model.SysResource, 0)

	err := global.DB.Model(resource).Association(model.SYSSUBRESOURCES).Find(&oldResources)
	if err != nil {
		return err
	}

	resources := append(oldResources, resource.SysSubResources...)
	return global.DB.Model(resource).Association(model.SYSSUBRESOURCES).Replace(resources)
}

func ClearSubResources(resource *model.SysResource) error {
	return global.DB.Model(resource).Association(model.SYSSUBRESOURCES).Delete(resource.SysSubResources)
}

func CreateSysResource(resource *model.SysResource) error {

	return global.DB.Create(resource).Error
}

func GetAllSysResource(t int) (model.SysResources, error) {
	list := model.SysResources{}
	db := global.DB.Model(&model.SysResource{})
	if t != 0 {
		db = db.Where("type = ?", t)
	}
	err := db.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil

}
