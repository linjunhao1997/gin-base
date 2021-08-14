package base

import (
	"fmt"
	"gin-base/component/db"
	"gorm.io/gorm"
)

type Range struct {
	start interface{}
	end   interface{}
}

type SearchParam struct {
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
	Eq       map[string]interface{} `json:"eq"`
	Like     map[string]string      `json:"like"`
	Range    map[string]*Range      `json:"range"`
	Sort     map[string]int         `json:"sort"`
}

type AllowField map[string]struct{}

type LoadField map[string]struct{}

func NewAllowField(fieldNames ...string) AllowField {
	a := make(AllowField, 0)
	for _, name := range fieldNames {
		a[name] = struct{}{}
	}
	return a
}

func NewLoadField(fieldNames ...string) LoadField {
	a := make(LoadField, 0)
	for _, name := range fieldNames {
		a[name] = struct{}{}
	}
	return a
}

func (a AllowField) contains(fieldName string) bool {
	if _, ok := a[fieldName]; ok {
		return true
	}
	return false
}

func (param *SearchParam) Validate(allowField AllowField) error {
	for k, _ := range param.Eq {
		if !allowField.contains(k) {
			return fmt.Errorf("不允许传入字段[%s]", k)
		}
	}
	for k, _ := range param.Like {
		if !allowField.contains(k) {
			return fmt.Errorf("不允许传入字段[%s]", k)
		}
	}
	for k, _ := range param.Range {
		if !allowField.contains(k) {
			return fmt.Errorf("不允许传入字段[%s]", k)
		}
	}
	for k, _ := range param.Sort {
		if !allowField.contains(k) {
			return fmt.Errorf("不允许传入字段[%s]", k)
		}
	}
	return nil
}

func (param *SearchParam) Search(loadField LoadField) *gorm.DB {
	db := db.RABC
	for field := range loadField {
		db = db.Preload(field)
	}
	db = db.Where(param.Eq)
	// 不支持"%"开头的like查询，效率太低
	for k, v := range param.Like {
		db = db.Where(fmt.Sprintf("%s LIKE ?", k), v+"%")
	}
	for k, v := range param.Range {
		if v.start != nil {
			db = db.Where(fmt.Sprintf("%s >= ?", k), v.start)
		}
		if v.end != nil {
			db = db.Where(fmt.Sprintf("%s < ?", k), v.end)
		}
	}
	for k, v := range param.Sort {
		if v == 1 {
			db = db.Order(k)
		}
		if v == -1 {
			db = db.Order(k + " desc")
		}
	}
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}

	offset := (param.Page - 1) * param.PageSize
	return db.Offset(offset).Limit(param.PageSize)
}
