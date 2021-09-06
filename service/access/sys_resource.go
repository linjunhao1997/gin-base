package service

import (
	"gin-base/global"
	model "gin-base/model/access"
	"gorm.io/gorm"
)

func SearchSysResources(db *gorm.DB) ([]model.SysResource, error) {

	resources := make([]model.SysResource, 0)
	if err := db.Find(&resources).Error; err != nil {
		return nil, err
	}
	return resources, nil
}

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
