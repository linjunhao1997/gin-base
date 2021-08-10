package handler

import "net/http"

// 错误处理的结构体
type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
}

var (
	Success      = NewError(http.StatusOK, 0, "success")
	ValidIdError = NewError(http.StatusBadRequest, 1, "id非法")
	BadRequest   = NewError(http.StatusBadRequest, 200400, "")
	ServerError  = NewError(http.StatusInternalServerError, 200500, "系统异常，请稍后重试!")
	NotFound     = NewError(http.StatusNotFound, 200404, http.StatusText(http.StatusNotFound))
)

func (e *Error) Error() string {
	return e.Msg
}

func NewError(statusCode, Code int, msg string) *Error {
	return &Error{
		StatusCode: statusCode,
		Code:       Code,
		Msg:        msg,
	}
}

func NewJsonError(msg string) *Error {
	return &Error{
		StatusCode: http.StatusBadRequest,
		Code:       1,
		Msg:        msg,
	}
}
