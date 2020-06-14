package responses

// SuccessResponse ...
type SuccessResponse struct {
}

// RespCode ...
type RespCode struct {
	Code       int    `json:"code"`
	Msg        string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (code *RespCode) Error() string {
	return code.Msg
}

// WithStatus ...
func (code *RespCode) WithStatus(status int) *RespCode {
	code.HTTPStatus = status
	return code
}

// NewErrorCode ...
func NewErrorCode(code int, msg string) *RespCode {
	return &RespCode{
		Code: code,
		Msg:  msg,
	}
}
