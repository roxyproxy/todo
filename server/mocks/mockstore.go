// Code generated by MockGen. DO NOT EDIT.
// Source: todo/storage (interfaces: Storage)

// Package mockstore is a generated GoMock package.
package mockstore

import (
	reflect "reflect"
	model "todo/model"
	storage "todo/storage"

	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddItem mocks base method.
func (m *MockStorage) AddItem(arg0 model.TodoItem) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddItem", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddItem indicates an expected call of AddItem.
func (mr *MockStorageMockRecorder) AddItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddItem", reflect.TypeOf((*MockStorage)(nil).AddItem), arg0)
}

// AddUser mocks base method.
func (m *MockStorage) AddUser(arg0 model.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStorageMockRecorder) AddUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStorage)(nil).AddUser), arg0)
}

// DeleteItem mocks base method.
func (m *MockStorage) DeleteItem(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteItem", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteItem indicates an expected call of DeleteItem.
func (mr *MockStorageMockRecorder) DeleteItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteItem", reflect.TypeOf((*MockStorage)(nil).DeleteItem), arg0)
}

// DeleteUser mocks base method.
func (m *MockStorage) DeleteUser(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockStorageMockRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockStorage)(nil).DeleteUser), arg0)
}

// GetAllItems mocks base method.
func (m *MockStorage) GetAllItems(arg0 storage.TodoFilter) ([]model.TodoItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllItems", arg0)
	ret0, _ := ret[0].([]model.TodoItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllItems indicates an expected call of GetAllItems.
func (mr *MockStorageMockRecorder) GetAllItems(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllItems", reflect.TypeOf((*MockStorage)(nil).GetAllItems), arg0)
}

// GetAllUsers mocks base method.
func (m *MockStorage) GetAllUsers(arg0 storage.UserFilter) ([]model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUsers", arg0)
	ret0, _ := ret[0].([]model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsers indicates an expected call of GetAllUsers.
func (mr *MockStorageMockRecorder) GetAllUsers(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsers", reflect.TypeOf((*MockStorage)(nil).GetAllUsers), arg0)
}

// GetItem mocks base method.
func (m *MockStorage) GetItem(arg0 string) (model.TodoItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetItem", arg0)
	ret0, _ := ret[0].(model.TodoItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockStorageMockRecorder) GetItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockStorage)(nil).GetItem), arg0)
}

// GetUser mocks base method.
func (m *MockStorage) GetUser(arg0 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStorageMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStorage)(nil).GetUser), arg0)
}

// UpdateItem mocks base method.
func (m *MockStorage) UpdateItem(arg0 model.TodoItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateItem", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateItem indicates an expected call of UpdateItem.
func (mr *MockStorageMockRecorder) UpdateItem(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateItem", reflect.TypeOf((*MockStorage)(nil).UpdateItem), arg0)
}

// UpdateUser mocks base method.
func (m *MockStorage) UpdateUser(arg0 model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockStorageMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockStorage)(nil).UpdateUser), arg0)
}
