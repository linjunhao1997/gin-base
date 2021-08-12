package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Err  string      `json:"error"`
	Ok   bool        `json:"success"`
	Data interface{} `json:"data"`
}

func (g *Gin) RespSuccess(data interface{}, msg string) {
	if msg == "" {
		msg = GetMsg(SUCCESS)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Ok:   true,
	})
}

func (g *Gin) RespServiceFail(failCode int, msg string) {
	if msg == "" {
		msg = GetMsg(Fail)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: failCode,
		Msg:  msg,
		Ok:   false,
	})
}

func (g *Gin) RespForbidden(msg string) {
	if msg == "" {
		msg = GetMsg(FORBIDDEN)
	}
	g.C.JSON(http.StatusForbidden, Response{
		Code: FORBIDDEN,
		Msg:  msg,
		Ok:   false,
	})
}

func (g *Gin) RespUnauthorized(msg string) {
	if msg == "" {
		msg = GetMsg(UNAUTHORIZED)
	}
	g.C.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Msg:  msg,
		Ok:   false,
	})
}

func (g *Gin) RespError(err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: ERROR,
		Msg:  msg,
		Err:  err.Error(),
		Ok:   false,
	})
}

func (g *Gin) RespNewError(httpCode, errCode int, err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  msg,
		Err:  err.Error(),
		Ok:   false,
	})
}
