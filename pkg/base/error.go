package base

var MsgFlags = map[int]string{
	SUCCESS:        "请求成功",
	ERROR:          "请求失败",
	INVALID_PARAMS: "请求参数错误",
	FORBIDDEN:      "禁止访问",
	UNAUTHORIZED:   "认证失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
