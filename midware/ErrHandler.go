package midware

import (
	"errors"
	"gin-base/handler"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		c.Next()

		if length := len(c.Errors); length > 0 {
			g := handler.Gin{C: c}
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					if len(errs) > 0 {
						err := errs[0]
						g.RespNewError(http.StatusBadRequest, handler.INVALID_PARAMS, errors.New(err.Translate(handler.Trans)), "")
						return
					}
				} else {
					g.RespNewError(http.StatusBadRequest, handler.INVALID_PARAMS, err, "")
					return
				}
			}

		}
	}
}
