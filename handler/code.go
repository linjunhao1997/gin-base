package handler

import "net/http"

const (
	SUCCESS        = http.StatusOK
	ERROR          = http.StatusInternalServerError
	INVALID_PARAMS = http.StatusBadRequest
	UNAUTHORIZED   = http.StatusUnauthorized
	FORBIDDEN      = http.StatusForbidden

	// 业务失败码
	Fail = 10001
)
