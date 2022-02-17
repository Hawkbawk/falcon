// Code generated by MockGen. DO NOT EDIT.
// Source: io (interfaces: ReadCloser)

// Package mock_io is a generated GoMock package.
package mock_io

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReadCloser is a mock of ReadCloser interface.
type MockReadCloser struct {
	ctrl     *gomock.Controller
	recorder *MockReadCloserMockRecorder
}

// MockReadCloserMockRecorder is the mock recorder for MockReadCloser.
type MockReadCloserMockRecorder struct {
	mock *MockReadCloser
}

// NewMockReadCloser creates a new mock instance.
func NewMockReadCloser(ctrl *gomock.Controller) *MockReadCloser {
	mock := &MockReadCloser{ctrl: ctrl}
	mock.recorder = &MockReadCloserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReadCloser) EXPECT() *MockReadCloserMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockReadCloser) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReadCloserMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReadCloser)(nil).Close))
}

// Read mocks base method.
func (m *MockReadCloser) Read(arg0 []byte) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockReadCloserMockRecorder) Read(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockReadCloser)(nil).Read), arg0)
}
