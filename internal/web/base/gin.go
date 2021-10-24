package base

import (
	model "gin-base/internal/model/access"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) EnsureSysUser() *model.SysUser {
	user, ok := g.C.Get("userInfo")
	if !ok {
		return nil
	}
	sysUser := user.(*model.SysUser)
	return sysUser
}

func EnsureSysUser(c gin.Context) *model.SysUser {
	user, _ := c.Get("userInfo")
	sysUser := user.(*model.SysUser)
	return sysUser
}

func (g *Gin) RespSuccess(data interface{}, msg string) {
	if msg == "" {
		msg = GetMsg(SUCCESS)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Ok:   true,
		Msg:  msg,
	})
}

func (g *Gin) RespServiceFail(failCode int, msg string) {
	if msg == "" {
		msg = GetMsg(Fail)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: failCode,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespForbidden(msg string) {
	if msg == "" {
		msg = GetMsg(FORBIDDEN)
	}
	g.C.Abort()
	g.C.JSON(http.StatusForbidden, Response{
		Code: FORBIDDEN,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespUnauthorized(msg string) {
	if msg == "" {
		msg = GetMsg(UNAUTHORIZED)
	}
	g.C.Abort()
	g.C.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespError(err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.Abort()
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: ERROR,
		Err:  err.Error(),
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespNewError(httpCode, errCode int, err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.Abort()
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Err:  err.Error(),
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) Abort(err error) {
	g.C.Abort()
	g.RespError(err, "")
}

type Controller struct {
}

func (controller *Controller) Wrap(f func(g *Gin)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(&Gin{C: c})
	}
}
