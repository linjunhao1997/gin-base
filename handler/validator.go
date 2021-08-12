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
	Trans ut.Translator
)

func validate(data interface{}) error {
	err := v.Struct(data)
	return err
}

func ValidateId(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return id, false
	}
	return id, true
}

func ValidateJson(c *gin.Context, body interface{}) bool {
	err := c.ShouldBind(body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return false
	}

	err = validate(body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return false
	}

	return true
}

func init() {
	uni := ut.New(zh.New(), zh.New())
	Trans, _ = uni.GetTranslator("zh")

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	zh_translations.RegisterDefaultTranslations(v, Trans)
}
