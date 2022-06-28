package common

import "encoding/json"

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func newErr(code int, message string) func(err error) *Err {
	return func(err error) *Err {
		var detail string
		if err != nil {
			detail = err.Error()
		}

		return &Err{
			Code:    code,
			Message: message,
			Detail:  detail,
		}
	}
}

func (e Err) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

var (
	ErrInternal           = newErr(1, "please send this message to developer")
	ErrInvalidRequestBody = newErr(2, "cannot parse request body")

	ErrNotFoundBook = newErr(100, "not found book")
	ErrNotFoundUser = newErr(101, "not found user")

	ErrStartRent = newErr(200, "cannot start rent")
	ErrEndRent   = newErr(201, "cannot end rent")
)
