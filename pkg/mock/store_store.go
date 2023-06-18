// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	store "github.com/raahii/haraiai/pkg/store"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreatePayment mocks base method.
func (m *MockStore) CreatePayment(arg0 string, arg1 *store.Payment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockStoreMockRecorder) CreatePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockStore)(nil).CreatePayment), arg0, arg1)
}

// DeleteGroup mocks base method.
func (m *MockStore) DeleteGroup(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockStoreMockRecorder) DeleteGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockStore)(nil).DeleteGroup), arg0)
}

// GetGroup mocks base method.
func (m *MockStore) GetGroup(arg0 string) (*store.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", arg0)
	ret0, _ := ret[0].(*store.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockStoreMockRecorder) GetGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockStore)(nil).GetGroup), arg0)
}

// SaveGroup mocks base method.
func (m *MockStore) SaveGroup(arg0 *store.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveGroup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveGroup indicates an expected call of SaveGroup.
func (mr *MockStoreMockRecorder) SaveGroup(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveGroup", reflect.TypeOf((*MockStore)(nil).SaveGroup), arg0)
}
