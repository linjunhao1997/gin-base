package base

import (
	"fmt"
	"gin-base/global/db"
	"gin-base/model"
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

type Pagination struct {
	List     interface{} `json:"list"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
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

func (param *SearchParam) Search(field ...string) *gorm.DB {

	fieldName := NewLoadField(field...)

	db := db.DB
	for k := range fieldName {
		db = db.Preload(k)
	}
	db = db.Scopes(EqFunc(param), LikeFunc(param), RangeFunc(param), SortFunc(param), PaginateFunc(param))
	return db
}

func (param *SearchParam) CountTotal(model interface{}) int {
	db := db.DB
	var total int64
	db.Model(model).Scopes(EqFunc(param), LikeFunc(param), RangeFunc(param)).Count(&total)
	return int(total)
}

func (param *SearchParam) NewPagination(data model.Models) Pagination {
	pagination := Pagination{
		Page:     param.Page,
		PageSize: param.PageSize,
		Total:    param.CountTotal(data.GetModel()),
	}

	if data.GetSize() == 0 {
		pagination.List = make([]interface{}, 0)
	}

	pagination.List = data
	return pagination
}

func EqFunc(param *SearchParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(param.Eq)
		return db
	}
}

func LikeFunc(param *SearchParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 不支持"%"开头的like查询，效率太低
		for k, v := range param.Like {
			db = db.Where(fmt.Sprintf("%s LIKE ?", k), v+"%")
		}
		return db
	}
}

func RangeFunc(param *SearchParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range param.Range {
			if v.start != nil {
				db = db.Where(fmt.Sprintf("%s >= ?", k), v.start)
			}
			if v.end != nil {
				db = db.Where(fmt.Sprintf("%s < ?", k), v.end)
			}
		}
		return db
	}
}

func SortFunc(param *SearchParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range param.Sort {
			if v == 1 {
				db = db.Order(k)
			}
			if v == -1 {
				db = db.Order(k + " desc")
			}
		}
		return db
	}
}

func PaginateFunc(param *SearchParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if param.Page == 0 {
			param.Page = 1
		}
		if param.PageSize == 0 {
			param.PageSize = 10
		}

		offset := (param.Page - 1) * param.PageSize
		return db.Offset(offset).Limit(param.PageSize)
	}
}
