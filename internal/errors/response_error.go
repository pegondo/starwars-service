package errors

import "fmt"

// ResponseError is the error structure used in the responses.
type ResponseError struct {
	// ErrorCode is the error code.
	ErrorCode string `json:"error_code"`
	// ErrorMessage is the error message.
	ErrorMessage string `json:"error_message"`
}

func (err *ResponseError) Error() string {
	return fmt.Sprintf("%s - %s", err.ErrorCode, err.ErrorMessage)
}

// New creates and returns a ResponseError error.
func New(errCode, errMsg string) error {
	return &ResponseError{
		ErrorCode:    errCode,
		ErrorMessage: errMsg,
	}
}

// InternalServerError returns a ResponseError for internal server errors.
func InternalServerError() error {
	return New(InternalServerErrorCode, InternalServerErrorMsg)
}
