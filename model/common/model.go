package common

type Model struct {
	ID int `gorm:"column:id;primary_key" json:"id"`
}
