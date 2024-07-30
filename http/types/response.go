package types

type Code int

const (
	CodeOk Code = iota
	CodeGeneralError
)

type Response struct {
	// RequestId string `json:"requestId"`

	// Success bool   `json:"success"`
	Code Code   `json:"code"`
	Msg  string `json:"msg"`

	Data any `json:"data"`
}
