// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/chatbotgang/go-clean-architecture-template/internal/app/service/auth (interfaces: AuthServer)

// Package automock is a generated GoMock package.
package automock

import (
	context "context"
	reflect "reflect"

	common "github.com/chatbotgang/go-clean-architecture-template/internal/domain/common"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthServer is a mock of AuthServer interface.
type MockAuthServer struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServerMockRecorder
}

// MockAuthServerMockRecorder is the mock recorder for MockAuthServer.
type MockAuthServerMockRecorder struct {
	mock *MockAuthServer
}

// NewMockAuthServer creates a new mock instance.
func NewMockAuthServer(ctrl *gomock.Controller) *MockAuthServer {
	mock := &MockAuthServer{ctrl: ctrl}
	mock.recorder = &MockAuthServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthServer) EXPECT() *MockAuthServerMockRecorder {
	return m.recorder
}

// AuthenticateAccount mocks base method.
func (m *MockAuthServer) AuthenticateAccount(arg0 context.Context, arg1, arg2 string) common.Error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthenticateAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(common.Error)
	return ret0
}

// AuthenticateAccount indicates an expected call of AuthenticateAccount.
func (mr *MockAuthServerMockRecorder) AuthenticateAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthenticateAccount", reflect.TypeOf((*MockAuthServer)(nil).AuthenticateAccount), arg0, arg1, arg2)
}

// RegisterAccount mocks base method.
func (m *MockAuthServer) RegisterAccount(arg0 context.Context, arg1, arg2 string) (string, common.Error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(common.Error)
	return ret0, ret1
}

// RegisterAccount indicates an expected call of RegisterAccount.
func (mr *MockAuthServerMockRecorder) RegisterAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterAccount", reflect.TypeOf((*MockAuthServer)(nil).RegisterAccount), arg0, arg1, arg2)
}