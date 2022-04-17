package accessmodel

type SysMenu struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Title     string      `gorm:"column:title" json:"title"`
	Icon      string      `gorm:"column:icon" json:"icon"`
	Url       string      `gorm:"column:url" json:"url"`
	ParentId  int         `gorm:"column:parent_id" json:"parent"`
	Desc      string      `gorm:"column:description" json:"desc"`
	Sorts     int         `gorm:"column:sort" json:"sorts"`
	Disabled  int         `gorm:"column:disabled" json:"disabled"`
	SysPowers []*SysPower `json:"powers"`
}

func (menu *SysMenu) GetResourceName() string {
	return "sysMenus"
}

type SysMenuBody struct {
	Title    *string `json:"title"`
	Icon     *string `json:"icon"`
	Url      *string `json:"url"`
	ParentId *int    `json:"parent"`
	Desc     *string `json:"desc"`
	Sorts    *int    `json:"sorts"`
	Disabled *int    `json:"disabled"`
}
