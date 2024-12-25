// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go
//
// Generated by this command:
//
//	mockgen -source=logger.go -destination=mocks/logger.go
//

// Package mock_logger is a generated GoMock package.
package mock_logger

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	zap "go.uber.org/zap"
)

// MockLoggerInterface is a mock of LoggerInterface interface.
type MockLoggerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerInterfaceMockRecorder
	isgomock struct{}
}

// MockLoggerInterfaceMockRecorder is the mock recorder for MockLoggerInterface.
type MockLoggerInterfaceMockRecorder struct {
	mock *MockLoggerInterface
}

// NewMockLoggerInterface creates a new mock instance.
func NewMockLoggerInterface(ctrl *gomock.Controller) *MockLoggerInterface {
	mock := &MockLoggerInterface{ctrl: ctrl}
	mock.recorder = &MockLoggerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerInterface) EXPECT() *MockLoggerInterfaceMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLoggerInterface) Debug(message string, fields ...zap.Field) {
	m.ctrl.T.Helper()
	varargs := []any{message}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerInterfaceMockRecorder) Debug(message any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{message}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLoggerInterface)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLoggerInterface) Error(message string, fields ...zap.Field) {
	m.ctrl.T.Helper()
	varargs := []any{message}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerInterfaceMockRecorder) Error(message any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{message}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLoggerInterface)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *MockLoggerInterface) Fatal(message string, fields ...zap.Field) {
	m.ctrl.T.Helper()
	varargs := []any{message}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockLoggerInterfaceMockRecorder) Fatal(message any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{message}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockLoggerInterface)(nil).Fatal), varargs...)
}

// Info mocks base method.
func (m *MockLoggerInterface) Info(message string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", message)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerInterfaceMockRecorder) Info(message any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLoggerInterface)(nil).Info), message)
}
