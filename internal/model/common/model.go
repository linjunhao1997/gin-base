package model

import "github.com/jinzhu/copier"

type Model interface {
	GetResourceName() string
}

func Mapping(dst interface{}, src interface{}) error {
	return copier.CopyWithOption(dst, src, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
}
