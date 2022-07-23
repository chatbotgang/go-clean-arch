package common

import (
	"net/http"
	"strconv"
	"strings"
)

type Error interface {
	Error() string
	// ClientMsg returns messages clients should know
	ClientMsg() string
}

type ErrorOption func(Error)

func WithMsg(msg string) ErrorOption {
	return func(e Error) {
		switch err := e.(type) {
		case *BaseError:
			err.clientMsg = msg
		}
	}
}

func WithStatus(status int) ErrorOption {
	return func(e Error) {
		switch err := e.(type) {
		case *BaseError:
			err.remoteStatus = status
		}
	}
}

func WithDetail(detail map[string]interface{}) ErrorOption {
	return func(e Error) {
		switch err := e.(type) {
		case *BaseError:
			err.detail = detail
		}
	}
}

type ErrorCode struct {
	Name       string
	StatusCode int
	Category   string // backward-compatible the current error format TODO: remove after deprecation
}

/*
	General error codes
*/

var ErrorCodeInternalProcess = ErrorCode{
	Name:       "INTERNAL_PROCESS",
	StatusCode: http.StatusInternalServerError,
	Category:   "INTERNAL_ERROR",
}

/*
	Authentication and Authorization error codes
*/

var ErrorCodeAuthPermissionDenied = ErrorCode{
	Name:       "AUTH_PERMISSION_DENIED",
	StatusCode: http.StatusForbidden,
	Category:   "RESOURCE_ERROR",
}

var ErrorCodeAuthNotAuthenticated = ErrorCode{
	Name:       "AUTH_NOT_AUTHENTICATED",
	StatusCode: http.StatusUnauthorized,
	Category:   "PARAMETER_ERROR",
}

/*
	Resource-related error codes
*/

var ErrorCodeResourceNotFound = ErrorCode{
	Name:       "RESOURCE_NOT_FOUND",
	StatusCode: http.StatusNotFound,
	Category:   "RESOURCE_ERROR",
}

/*
	Parameter-related error codes
*/

var ErrorCodeParameterInvalid = ErrorCode{
	Name:       "PARAMETER_INVALID",
	StatusCode: http.StatusBadRequest,
	Category:   "PARAMETER_ERROR",
}

/*
	Remote server-related error codes
*/

var ErrorCodeRemoteProcess = ErrorCode{
	Name:       "REMOTE_PROCESS_ERROR",
	StatusCode: http.StatusBadGateway,
	Category:   "REMOTE_SERVER_ERROR",
}

// BaseError used for expressing errors occurring in application.
type BaseError struct {
	err          error
	code         ErrorCode
	clientMsg    string
	remoteStatus int // proxy HTTP status code
	detail       map[string]interface{}
}

func NewError(code ErrorCode, err error, opts ...ErrorOption) Error {
	if err, ok := err.(Error); ok {
		return err
	}

	e := BaseError{code: code, err: err}
	for _, o := range opts {
		o(&e)
	}
	return e
}

func (e BaseError) Error() string {
	var msgs []string
	if e.remoteStatus != 0 {
		msgs = append(msgs, strconv.Itoa(e.remoteStatus))
	}
	if e.err != nil {
		msgs = append(msgs, e.err.Error())
	}
	if e.clientMsg != "" {
		msgs = append(msgs, e.clientMsg)
	}

	return strings.Join(msgs, ": ")
}

func (e BaseError) Name() string {
	if e.code.Name == "" {
		return "UNKNOWN_ERROR"
	}
	return e.code.Name
}

func (e BaseError) Category() string {
	if e.code.Category == "" {
		return "UNKNOWN_ERROR"
	}
	return e.code.Category
}

func (e BaseError) ClientMsg() string {
	return e.clientMsg
}

func (e BaseError) HTTPStatus() int {
	if e.code.StatusCode == 0 {
		return http.StatusInternalServerError
	}
	return e.code.StatusCode
}

func (e BaseError) RemoteHTTPStatus() int {
	return e.remoteStatus
}

func (e BaseError) Detail() map[string]interface{} {
	return e.detail
}
