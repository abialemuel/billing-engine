// Code generated by MockGen. DO NOT EDIT.
// Source: config/config.go
//
// Generated by this command:
//
//	mockgen -source=config/config.go -destination=mocks/config/config.go -package=config Config
//

// Package config is a generated GoMock package.
package config

import (
	reflect "reflect"

	config "github.com/abialemuel/billing-engine/config"
	gomock "go.uber.org/mock/gomock"
)

// MockConfig is a mock of Config interface.
type MockConfig struct {
	ctrl     *gomock.Controller
	recorder *MockConfigMockRecorder
}

// MockConfigMockRecorder is the mock recorder for MockConfig.
type MockConfigMockRecorder struct {
	mock *MockConfig
}

// NewMockConfig creates a new mock instance.
func NewMockConfig(ctrl *gomock.Controller) *MockConfig {
	mock := &MockConfig{ctrl: ctrl}
	mock.recorder = &MockConfigMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfig) EXPECT() *MockConfigMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConfig) Get() *config.MainConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*config.MainConfig)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockConfigMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfig)(nil).Get))
}

// Init mocks base method.
func (m *MockConfig) Init(configPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", configPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockConfigMockRecorder) Init(configPath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockConfig)(nil).Init), configPath)
}
