package model

type SysUser struct {
	Id       int    `gorm:"column:id;primary_key" json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json："password"`
}
