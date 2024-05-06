// Code generated by mockery v2.43.0. DO NOT EDIT.

package common

import mock "github.com/stretchr/testify/mock"

// MockExecutorProvider is an autogenerated mock type for the ExecutorProvider type
type MockExecutorProvider struct {
	mock.Mock
}

// Acquire provides a mock function with given fields: config
func (_m *MockExecutorProvider) Acquire(config *RunnerConfig) (ExecutorData, error) {
	ret := _m.Called(config)

	if len(ret) == 0 {
		panic("no return value specified for Acquire")
	}

	var r0 ExecutorData
	var r1 error
	if rf, ok := ret.Get(0).(func(*RunnerConfig) (ExecutorData, error)); ok {
		return rf(config)
	}
	if rf, ok := ret.Get(0).(func(*RunnerConfig) ExecutorData); ok {
		r0 = rf(config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ExecutorData)
		}
	}

	if rf, ok := ret.Get(1).(func(*RunnerConfig) error); ok {
		r1 = rf(config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CanCreate provides a mock function with given fields:
func (_m *MockExecutorProvider) CanCreate() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CanCreate")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Create provides a mock function with given fields:
func (_m *MockExecutorProvider) Create() Executor {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 Executor
	if rf, ok := ret.Get(0).(func() Executor); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(Executor)
		}
	}

	return r0
}

// GetConfigInfo provides a mock function with given fields: input, output
func (_m *MockExecutorProvider) GetConfigInfo(input *RunnerConfig, output *ConfigInfo) {
	_m.Called(input, output)
}

// GetDefaultShell provides a mock function with given fields:
func (_m *MockExecutorProvider) GetDefaultShell() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetDefaultShell")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// GetFeatures provides a mock function with given fields: features
func (_m *MockExecutorProvider) GetFeatures(features *FeaturesInfo) error {
	ret := _m.Called(features)

	if len(ret) == 0 {
		panic("no return value specified for GetFeatures")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*FeaturesInfo) error); ok {
		r0 = rf(features)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Release provides a mock function with given fields: config, data
func (_m *MockExecutorProvider) Release(config *RunnerConfig, data ExecutorData) {
	_m.Called(config, data)
}

// NewMockExecutorProvider creates a new instance of MockExecutorProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExecutorProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExecutorProvider {
	mock := &MockExecutorProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
