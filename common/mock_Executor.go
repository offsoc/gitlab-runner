// Code generated by mockery v2.14.0. DO NOT EDIT.

package common

import mock "github.com/stretchr/testify/mock"

// MockExecutor is an autogenerated mock type for the Executor type
type MockExecutor struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields:
func (_m *MockExecutor) Cleanup() {
	_m.Called()
}

// Finish provides a mock function with given fields: err
func (_m *MockExecutor) Finish(err error) {
	_m.Called(err)
}

// GetCurrentStage provides a mock function with given fields:
func (_m *MockExecutor) GetCurrentStage() ExecutorStage {
	ret := _m.Called()

	var r0 ExecutorStage
	if rf, ok := ret.Get(0).(func() ExecutorStage); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ExecutorStage)
	}

	return r0
}

// Prepare provides a mock function with given fields: options
func (_m *MockExecutor) Prepare(options ExecutorPrepareOptions) error {
	ret := _m.Called(options)

	var r0 error
	if rf, ok := ret.Get(0).(func(ExecutorPrepareOptions) error); ok {
		r0 = rf(options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Run provides a mock function with given fields: cmd
func (_m *MockExecutor) Run(cmd ExecutorCommand) error {
	ret := _m.Called(cmd)

	var r0 error
	if rf, ok := ret.Get(0).(func(ExecutorCommand) error); ok {
		r0 = rf(cmd)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetCurrentStage provides a mock function with given fields: stage
func (_m *MockExecutor) SetCurrentStage(stage ExecutorStage) {
	_m.Called(stage)
}

// Shell provides a mock function with given fields:
func (_m *MockExecutor) Shell() *ShellScriptInfo {
	ret := _m.Called()

	var r0 *ShellScriptInfo
	if rf, ok := ret.Get(0).(func() *ShellScriptInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ShellScriptInfo)
		}
	}

	return r0
}

type mockConstructorTestingTNewMockExecutor interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockExecutor creates a new instance of MockExecutor. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockExecutor(t mockConstructorTestingTNewMockExecutor) *MockExecutor {
	mock := &MockExecutor{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
