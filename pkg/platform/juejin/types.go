package juejin

type APIError struct {
	ErrMsg string `json:"err_msg"`
	ErrNo  int    `json:"err_no"`
}
