package accessmodel

type SysApi struct {
	ID     int    `gorm:"column:id;primary_key" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Desc   string `gorm:"column:description" json:"desc"`
	Url    string `gorm:"column:url" json:"url"`
	Method string `gorm:"column:method" json:"method"`
	Enable int    `gorm:"column:enable" json:"conditions"`
}

func (api *SysApi) GetResourceName() string {
	return "sysApis"
}

type SysApis []*SysApi

type SysApiBody struct {
	Name   *string `json:"name"`
	Desc   *string `json:"desc"`
	Url    *string `json:"url"`
	Method *string `json:"method" validate:"omitempty,oneof=POST GET"`
	Enable *int    `json:"conditions"`
}
