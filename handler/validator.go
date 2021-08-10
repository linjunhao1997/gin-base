package handler

import (
	"github.com/gin-gonic/gin"
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

func validate(data interface{}) (bool, string) {
	err := v.Struct(data)
	errs := err.(validator.ValidationErrors)
	if len(errs) > 0 {
		err := errs[0]
		return false, err.Translate(trans)
	}
	return true, ""
}

func ValidateId(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, ValidIdError)
		return id, false
	}
	return id, true
}

func ValidateJson(c *gin.Context, body interface{}) bool {
	err := c.ShouldBind(body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, NewJsonError(err.Error()))
		return false
	}

	ok, errStr := validate(body)
	if !ok {
		c.AbortWithError(http.StatusBadRequest, NewJsonError(errStr))
		return false
	}

	return true
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
