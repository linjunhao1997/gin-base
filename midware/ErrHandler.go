package midware

import (
	"gin-base/handler"
	"github.com/gin-gonic/gin"
)

func ErrHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		c.Next()

		if length := len(c.Errors); length > 0 {
			e := c.Errors[length-1]
			err := e.Err
			if err != nil {
				var Err *handler.Error
				if e, ok := err.(*handler.Error); ok {
					Err = e
				} else if e, ok := err.(error); ok {
					Err = handler.NewError(500, -1, e.Error())
				} else {
					Err = handler.ServerError
				}

				c.JSON(Err.StatusCode, Err)
				return
			}
		}

	}
}
