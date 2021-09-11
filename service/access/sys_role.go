package service

import (
	"gin-base/global"
	model "gin-base/model/access"
)

func GetApiResources(roleId ...int) ([]model.SysResource, error) {

	apis := make([]model.SysResource, 0)
	err := global.DB.Where("type = ? and disable = 0", model.API).Find(&apis).Error
	if err != nil {
		return apis, err
	}
	return apis, nil
}
func CreateSysRole(role *model.SysRole) error {
	return global.DB.Create(role).Error
}

func RelatedRoleResources(roleId int, resourceIds []int) error {
	role := model.SysRole{ID: roleId}
	resources := make([]model.SysResource, len(resourceIds))
	for i, resourceId := range resourceIds {
		resources[i] = model.SysResource{ID: resourceId}
	}
	return global.DB.Model(&role).Association("SysResources").Replace(&resources)
}
