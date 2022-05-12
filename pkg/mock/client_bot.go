// Code generated by MockGen. DO NOT EDIT.
// Source: bot.go

// Package mock is a generated GoMock package.
package mock

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	linebot "github.com/line/line-bot-sdk-go/v7/linebot"
)

// MockBotClient is a mock of BotClient interface.
type MockBotClient struct {
	ctrl     *gomock.Controller
	recorder *MockBotClientMockRecorder
}

// MockBotClientMockRecorder is the mock recorder for MockBotClient.
type MockBotClientMockRecorder struct {
	mock *MockBotClient
}

// NewMockBotClient creates a new mock instance.
func NewMockBotClient(ctrl *gomock.Controller) *MockBotClient {
	mock := &MockBotClient{ctrl: ctrl}
	mock.recorder = &MockBotClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBotClient) EXPECT() *MockBotClientMockRecorder {
	return m.recorder
}

// ParseRequest mocks base method.
func (m *MockBotClient) ParseRequest(arg0 *http.Request) ([]*linebot.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseRequest", arg0)
	ret0, _ := ret[0].([]*linebot.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseRequest indicates an expected call of ParseRequest.
func (mr *MockBotClientMockRecorder) ParseRequest(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseRequest", reflect.TypeOf((*MockBotClient)(nil).ParseRequest), arg0)
}

// ReplyMessage mocks base method.
func (m *MockBotClient) ReplyMessage(arg0 string, arg1 ...linebot.SendingMessage) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReplyMessage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplyMessage indicates an expected call of ReplyMessage.
func (mr *MockBotClientMockRecorder) ReplyMessage(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplyMessage", reflect.TypeOf((*MockBotClient)(nil).ReplyMessage), varargs...)
}

// ReplyTextMessage mocks base method.
func (m *MockBotClient) ReplyTextMessage(arg0 string, arg1 ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ReplyTextMessage", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplyTextMessage indicates an expected call of ReplyTextMessage.
func (mr *MockBotClientMockRecorder) ReplyTextMessage(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplyTextMessage", reflect.TypeOf((*MockBotClient)(nil).ReplyTextMessage), varargs...)
}
