package base

import (
	"errors"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	v     = validator.New()
	trans ut.Translator
)

func (g *Gin) validate(data interface{}) error {
	err := v.Struct(data)
	return err
}

func (g *Gin) ValidateId() (int, bool) {
	id, err := strconv.Atoi(g.C.Param("id"))
	if err != nil {
		g.RespNewError(http.StatusBadRequest, INVALID_PARAMS, err, "")
		return id, false
	}
	return id, true
}

func (g *Gin) ValidateJson(body interface{}) bool {
	err := g.C.ShouldBindJSON(body)
	if err != nil {
		g.RespNewError(http.StatusBadRequest, INVALID_PARAMS, err, "")
		return false
	}

	err = g.validate(body)
	if err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			if len(errs) > 0 {
				err := errs[0]
				g.RespNewError(http.StatusBadRequest, INVALID_PARAMS, errors.New(err.Translate(trans)), "")
			}
		} else {
			g.RespNewError(http.StatusBadRequest, INVALID_PARAMS, err, "")
		}
		return false
	}

	return true
}

func (g *Gin) ValidateAllowField(field AllowField) *SearchParam {

	var body SearchParam
	ok := g.ValidateJson(&body)
	if !ok {
		return nil
	}

	err := body.Validate(field)
	if err != nil {
		g.RespNewError(http.StatusBadRequest, INVALID_PARAMS, err, "")
		return nil
	}
	return &body
}

func init() {
	uni := ut.New(zh.New(), zh.New())
	trans, _ = uni.GetTranslator("zh")

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	zh_translations.RegisterDefaultTranslations(v, trans)
}
