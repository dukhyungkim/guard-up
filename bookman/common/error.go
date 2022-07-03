package common

import (
	"encoding/json"
	"net/http"
)

type Err struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Detail     string `json:"detail"`
	HTTPStatus int    `json:"-"`
}

func newErr(code int, message string, httpStatus int) func(err error) *Err {
	return func(err error) *Err {
		var detail string
		if err != nil {
			detail = err.Error()
		}

		return &Err{
			Code:       code,
			Message:    message,
			Detail:     detail,
			HTTPStatus: httpStatus,
		}
	}
}

func (e Err) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

var (
	ErrInternal           = newErr(1, "please send this message to developer", http.StatusInternalServerError)
	ErrInvalidRequestBody = newErr(2, "cannot parse request body", http.StatusBadRequest)
	ErrInvalidParam       = newErr(3, "cannot parse request param", http.StatusBadRequest)

	ErrNotFoundBook         = newErr(100, "not found book", http.StatusNotFound)
	ErrNotFoundUser         = newErr(101, "not found user", http.StatusNotFound)
	ErrNotFoundBookOrUser   = newErr(102, "not found book or user", http.StatusNotFound)
	ErrNotFoundRentalStatus = newErr(103, "not found rental status", http.StatusNoContent)

	ErrStartRent = newErr(200, "cannot start rent", http.StatusBadRequest)

	ErrNotFoundAction = newErr(300, "not found action", http.StatusBadRequest)
)
