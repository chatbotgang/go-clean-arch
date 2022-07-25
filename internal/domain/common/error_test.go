package common

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseError_Option(t *testing.T) {
	t.Parallel()

	msg := "random client message"
	status := http.StatusBadRequest
	detail := map[string]interface{}{
		"channel_id":  123,
		"member_name": "who am I?",
		"tag_id":      []int{1, 2, 3, 4},
	}

	// Test cases
	testCases := []struct {
		Name                string
		TestError           Error
		WithMsg             bool
		WitStatus           bool
		WitDeniedPermission bool
		WithDetail          bool
	}{
		{
			Name:      "with client msg",
			TestError: NewError(ErrorCodeInternalProcess, nil, WithMsg(msg)),
			WithMsg:   true,
		},
		{
			Name:      "with proxy HTTP status",
			TestError: NewError(ErrorCodeRemoteProcess, nil, WithStatus(status)),
			WitStatus: true,
		},
		{
			Name:       "with detail",
			TestError:  NewError(ErrorCodeInternalProcess, nil, WithDetail(detail)),
			WithDetail: true,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			err := c.TestError

			var baseError BaseError
			if errors.As(err, &baseError) {
				if c.WithMsg {
					assert.EqualValues(t, msg, baseError.ClientMsg())
				}
				if c.WitStatus {
					assert.EqualValues(t, status, baseError.RemoteHTTPStatus())
				}
				if c.WitDeniedPermission {
					assert.Contains(t, baseError.Error(), "no permission to")
				}
				if c.WithDetail {
					assert.Contains(t, baseError.Detail(), "channel_id")
					assert.Contains(t, baseError.Detail(), "member_name")
					assert.Equal(t, detail, baseError.Detail())
				}
			}
		})
	}
}

func TestBaseError_ErrorMapping(t *testing.T) {
	t.Parallel()

	// Test cases
	testCases := []struct {
		Name                   string
		TestError              Error
		ExpectErrorName        string
		ExpectHTTPStatus       int
		ExpectRemoteHTTPStatus int
	}{
		{
			Name:             "internal process",
			TestError:        NewError(ErrorCodeInternalProcess, nil),
			ExpectErrorName:  ErrorCodeInternalProcess.Name,
			ExpectHTTPStatus: http.StatusInternalServerError,
		},
		{
			Name:             "permission denied",
			TestError:        NewError(ErrorCodeAuthPermissionDenied, nil),
			ExpectErrorName:  ErrorCodeAuthPermissionDenied.Name,
			ExpectHTTPStatus: http.StatusForbidden,
		},
		{
			Name:             "not authenticated",
			TestError:        NewError(ErrorCodeAuthNotAuthenticated, nil),
			ExpectErrorName:  ErrorCodeAuthNotAuthenticated.Name,
			ExpectHTTPStatus: http.StatusUnauthorized,
		},
		{
			Name:             "invalid parameter",
			TestError:        NewError(ErrorCodeParameterInvalid, nil),
			ExpectErrorName:  ErrorCodeParameterInvalid.Name,
			ExpectHTTPStatus: http.StatusBadRequest,
		},
		{
			Name:             "resource not found",
			TestError:        NewError(ErrorCodeResourceNotFound, nil),
			ExpectErrorName:  ErrorCodeResourceNotFound.Name,
			ExpectHTTPStatus: http.StatusNotFound,
		},
		{
			Name:                   "remote process",
			TestError:              NewError(ErrorCodeRemoteProcess, nil, WithStatus(http.StatusBadRequest)),
			ExpectErrorName:        ErrorCodeRemoteProcess.Name,
			ExpectHTTPStatus:       http.StatusBadGateway,
			ExpectRemoteHTTPStatus: http.StatusBadRequest,
		},
		{
			Name:             "unknown error",
			TestError:        NewError(ErrorCode{}, nil),
			ExpectErrorName:  "UNKNOWN_ERROR",
			ExpectHTTPStatus: http.StatusInternalServerError,
		},
	}

	for i := range testCases {
		c := testCases[i]
		t.Run(c.Name, func(t *testing.T) {
			err := c.TestError

			var baseError BaseError
			if errors.As(err, &baseError) {
				assert.Equal(t, c.ExpectErrorName, baseError.Name())
				assert.Equal(t, c.ExpectHTTPStatus, baseError.HTTPStatus())
				assert.Equal(t, c.ExpectRemoteHTTPStatus, baseError.RemoteHTTPStatus())
			} else {
				t.Error("failed to match error type")
			}
		})
	}
}
