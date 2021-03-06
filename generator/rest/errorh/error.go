package errorh

import (
	"errors"
	"fmt"
	"net/http"
)

func NewInternalErrorf(code int, format string, args ...interface{}) *Error {
	return NewInternalError(code, fmt.Errorf(format, args...))
}

func NewInternalError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpErrorType = http.StatusInternalServerError
	return newError
}

func NewNotImplementedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotImplementedError(code, fmt.Errorf(format, args...))
}

func NewNotImplementedError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpErrorType = http.StatusNotImplemented
	return newError
}

func NewConflictErrorf(code int, format string, args ...interface{}) *Error {
	return NewConflictError(code, fmt.Errorf(format, args...))
}

func NewConflictError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpErrorType = http.StatusConflict
	return newError
}

func NewInvalidInputErrorf(code int, format string, args ...interface{}) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = fmt.Errorf(format, args...)
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpErrorType = http.StatusBadRequest
	return newError
}

func NewInvalidInputErrorSpecific(code int, fieldErrors []FieldError) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = errors.New("Input validation error")
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpErrorType = http.StatusBadRequest
	newError.FieldErrors = fieldErrors
	return newError
}

func NewNotFoundErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotFoundError(code, fmt.Errorf(format, args...))
}

func NewNotFoundError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = err.Error()
	newError.httpErrorType = http.StatusNotFound
	return newError
}

func NewNotAuthorizedErrorf(code int, format string, args ...interface{}) *Error {
	return NewNotAuthorizedError(code, fmt.Errorf(format, args...))
}

func NewNotAuthorizedError(code int, err error) *Error {
	newError := new(Error)
	newError.ErrorCode = code
	newError.underlyingError = err
	newError.ErrorMessage = newError.underlyingError.Error()
	newError.httpErrorType = http.StatusForbidden
	return newError
}

func (err Error) Error() string {
	return err.ErrorMessage
}

func (err Error) IsInternalError() bool {
	return err.httpErrorType == http.StatusInternalServerError
}

func (err Error) IsNotImplementedError() bool {
	return err.httpErrorType == http.StatusNotImplemented
}

func (err Error) IsInvalidInputError() bool {
	return err.httpErrorType == http.StatusBadRequest
}

func (err Error) GetFieldErrors() []FieldError {
	return err.FieldErrors
}

func (err Error) IsNotFoundError() bool {
	return err.httpErrorType == http.StatusNotFound
}

func (err Error) IsConflictError() bool {
	return err.httpErrorType == http.StatusConflict
}

func (err Error) IsNotAuthorizedError() bool {
	return err.httpErrorType == http.StatusForbidden
}

func (err Error) GetHttpCode() int {
	return err.httpErrorType
}

func (err Error) GetErrorCode() int {
	return err.ErrorCode
}

func GetFieldErrors(err error) []FieldError {
	return getFieldErrors(err)
}

func getFieldErrors(err error) []FieldError {
	if IsInvalidInputError(err) {
		e, _ := err.(InvalidInput)
		return e.GetFieldErrors()
	} else {
		return []FieldError{}
	}
}
