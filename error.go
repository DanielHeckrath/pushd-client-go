package client

import (
	"fmt"
)

const (
	UNEXPECTED_ERROR_MESSAGE = "Unexpected response from stats service: status=%d, body=%s"

	INVALID_PARAMETER_ERROR      = 40001
	ACCOUNT_NOT_EXISTS_ERROR     = 40002
	ACCOUNT_ALREADY_EXISTS_ERROR = 40003
)

type ClientError struct {
	Body string
	Code int
}

func (this *ClientError) Error() string {
	return this.Body
}

func IsError(err error, code int) bool {
	switch e := err.(type) {

	case *ClientError:
		return e.Code == code
	}

	return false
}

func newUnexpectedResponseError(code int, body string) error {
	return &ClientError{
		Body: fmt.Sprintf(UNEXPECTED_ERROR, code, body),
		Code: code,
	}
}
