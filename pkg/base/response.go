package base

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Err  string      `json:"error"`
	Ok   bool        `json:"success"`
	Msg  string      `json:"message"`
}
