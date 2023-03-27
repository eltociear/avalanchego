// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/avalanchego/snow/networking/handler (interfaces: Handler)

// Package handler is a generated GoMock package.
package handler

import (
	context "context"
	reflect "reflect"

	ids "github.com/ava-labs/avalanchego/ids"
	snow "github.com/ava-labs/avalanchego/snow"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockHandler) Context() *snow.ConsensusContext {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(*snow.ConsensusContext)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockHandlerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockHandler)(nil).Context))
}

// HealthCheck mocks base method.
func (m *MockHandler) HealthCheck(arg0 context.Context) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthCheck", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// HealthCheck indicates an expected call of HealthCheck.
func (mr *MockHandlerMockRecorder) HealthCheck(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthCheck", reflect.TypeOf((*MockHandler)(nil).HealthCheck), arg0)
}

// Len mocks base method.
func (m *MockHandler) Len() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Len")
	ret0, _ := ret[0].(int)
	return ret0
}

// Len indicates an expected call of Len.
func (mr *MockHandlerMockRecorder) Len() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Len", reflect.TypeOf((*MockHandler)(nil).Len))
}

// Push mocks base method.
func (m *MockHandler) Push(arg0 context.Context, arg1 Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Push", arg0, arg1)
}

// Push indicates an expected call of Push.
func (mr *MockHandlerMockRecorder) Push(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockHandler)(nil).Push), arg0, arg1)
}

// SetEngineManager mocks base method.
func (m *MockHandler) SetEngineManager(arg0 *EngineManager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetEngineManager", arg0)
}

// SetEngineManager indicates an expected call of SetEngineManager.
func (mr *MockHandlerMockRecorder) SetEngineManager(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetEngineManager", reflect.TypeOf((*MockHandler)(nil).SetEngineManager), arg0)
}

// SetOnStopped mocks base method.
func (m *MockHandler) SetOnStopped(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetOnStopped", arg0)
}

// SetOnStopped indicates an expected call of SetOnStopped.
func (mr *MockHandlerMockRecorder) SetOnStopped(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOnStopped", reflect.TypeOf((*MockHandler)(nil).SetOnStopped), arg0)
}

// ShouldHandle mocks base method.
func (m *MockHandler) ShouldHandle(arg0 ids.NodeID) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldHandle", arg0)
	ret0, _ := ret[0].(bool)
	return ret0
}

// ShouldHandle indicates an expected call of ShouldHandle.
func (mr *MockHandlerMockRecorder) ShouldHandle(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldHandle", reflect.TypeOf((*MockHandler)(nil).ShouldHandle), arg0)
}

// Start mocks base method.
func (m *MockHandler) Start(arg0 context.Context, arg1 bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", arg0, arg1)
}

// Start indicates an expected call of Start.
func (mr *MockHandlerMockRecorder) Start(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockHandler)(nil).Start), arg0, arg1)
}

// Stop mocks base method.
func (m *MockHandler) Stop(arg0 context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop", arg0)
}

// Stop indicates an expected call of Stop.
func (mr *MockHandlerMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockHandler)(nil).Stop), arg0)
}

// StopWithError mocks base method.
func (m *MockHandler) StopWithError(arg0 context.Context, arg1 error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopWithError", arg0, arg1)
}

// StopWithError indicates an expected call of StopWithError.
func (mr *MockHandlerMockRecorder) StopWithError(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopWithError", reflect.TypeOf((*MockHandler)(nil).StopWithError), arg0, arg1)
}

// Stopped mocks base method.
func (m *MockHandler) Stopped() chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stopped")
	ret0, _ := ret[0].(chan struct{})
	return ret0
}

// Stopped indicates an expected call of Stopped.
func (mr *MockHandlerMockRecorder) Stopped() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stopped", reflect.TypeOf((*MockHandler)(nil).Stopped))
}
