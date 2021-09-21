package service

import (
	"gin-base/global"
	model "gin-base/model/access"
	"strings"
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

func GetSysResources(roleId ...int) ([]*model.SysResource, error) {
	slice := make([]string, 4)
	slice[0] = model.SYSSUBRESOURCES
	for i := 1; i < 4; i++ {
		slice[i] = strings.Join([]string{slice[i-1], model.SYSSUBRESOURCES}, ".")
	}
	sysResource := make([]*model.SysResource, 0)

	if err := global.DB.Distinct().Table("sys_resource").Joins("left join sys_role_r_sys_resource on sys_resource.id = sys_role_r_sys_resource.sys_resource_id").Where("sys_role_r_sys_resource.sys_role_id IN ? and sys_resource.type = ?", roleId, model.MODULE).Preload(slice[3]).Find(&sysResource).Error; err != nil {
		return sysResource, err
	}
	return sysResource, nil
}
