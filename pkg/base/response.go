package base

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Err  string      `json:"error"`
	Ok   bool        `json:"success"`
	Data interface{} `json:"data"`
}
