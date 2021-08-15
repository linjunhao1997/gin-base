package model

type Base struct {
	ID int `gorm:"column:id;primary_key" json:"id"`
}

type Models interface {
	GetSize() int
	GetModel() interface{}
}
