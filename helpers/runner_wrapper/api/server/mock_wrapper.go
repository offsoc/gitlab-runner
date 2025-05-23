// Code generated by mockery v2.53.3. DO NOT EDIT.

package server

import (
	mock "github.com/stretchr/testify/mock"
	api "gitlab.com/gitlab-org/gitlab-runner/helpers/runner_wrapper/api"
)

// mockWrapper is an autogenerated mock type for the wrapper type
type mockWrapper struct {
	mock.Mock
}

type mockWrapper_Expecter struct {
	mock *mock.Mock
}

func (_m *mockWrapper) EXPECT() *mockWrapper_Expecter {
	return &mockWrapper_Expecter{mock: &_m.Mock}
}

// FailureReason provides a mock function with no fields
func (_m *mockWrapper) FailureReason() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for FailureReason")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockWrapper_FailureReason_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FailureReason'
type mockWrapper_FailureReason_Call struct {
	*mock.Call
}

// FailureReason is a helper method to define mock.On call
func (_e *mockWrapper_Expecter) FailureReason() *mockWrapper_FailureReason_Call {
	return &mockWrapper_FailureReason_Call{Call: _e.mock.On("FailureReason")}
}

func (_c *mockWrapper_FailureReason_Call) Run(run func()) *mockWrapper_FailureReason_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockWrapper_FailureReason_Call) Return(_a0 string) *mockWrapper_FailureReason_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockWrapper_FailureReason_Call) RunAndReturn(run func() string) *mockWrapper_FailureReason_Call {
	_c.Call.Return(run)
	return _c
}

// InitiateForcefulShutdown provides a mock function with no fields
func (_m *mockWrapper) InitiateForcefulShutdown() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for InitiateForcefulShutdown")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockWrapper_InitiateForcefulShutdown_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitiateForcefulShutdown'
type mockWrapper_InitiateForcefulShutdown_Call struct {
	*mock.Call
}

// InitiateForcefulShutdown is a helper method to define mock.On call
func (_e *mockWrapper_Expecter) InitiateForcefulShutdown() *mockWrapper_InitiateForcefulShutdown_Call {
	return &mockWrapper_InitiateForcefulShutdown_Call{Call: _e.mock.On("InitiateForcefulShutdown")}
}

func (_c *mockWrapper_InitiateForcefulShutdown_Call) Run(run func()) *mockWrapper_InitiateForcefulShutdown_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockWrapper_InitiateForcefulShutdown_Call) Return(_a0 error) *mockWrapper_InitiateForcefulShutdown_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockWrapper_InitiateForcefulShutdown_Call) RunAndReturn(run func() error) *mockWrapper_InitiateForcefulShutdown_Call {
	_c.Call.Return(run)
	return _c
}

// InitiateGracefulShutdown provides a mock function with given fields: req
func (_m *mockWrapper) InitiateGracefulShutdown(req api.InitGracefulShutdownRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for InitiateGracefulShutdown")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(api.InitGracefulShutdownRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockWrapper_InitiateGracefulShutdown_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InitiateGracefulShutdown'
type mockWrapper_InitiateGracefulShutdown_Call struct {
	*mock.Call
}

// InitiateGracefulShutdown is a helper method to define mock.On call
//   - req api.InitGracefulShutdownRequest
func (_e *mockWrapper_Expecter) InitiateGracefulShutdown(req interface{}) *mockWrapper_InitiateGracefulShutdown_Call {
	return &mockWrapper_InitiateGracefulShutdown_Call{Call: _e.mock.On("InitiateGracefulShutdown", req)}
}

func (_c *mockWrapper_InitiateGracefulShutdown_Call) Run(run func(req api.InitGracefulShutdownRequest)) *mockWrapper_InitiateGracefulShutdown_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(api.InitGracefulShutdownRequest))
	})
	return _c
}

func (_c *mockWrapper_InitiateGracefulShutdown_Call) Return(_a0 error) *mockWrapper_InitiateGracefulShutdown_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockWrapper_InitiateGracefulShutdown_Call) RunAndReturn(run func(api.InitGracefulShutdownRequest) error) *mockWrapper_InitiateGracefulShutdown_Call {
	_c.Call.Return(run)
	return _c
}

// Status provides a mock function with no fields
func (_m *mockWrapper) Status() api.Status {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Status")
	}

	var r0 api.Status
	if rf, ok := ret.Get(0).(func() api.Status); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(api.Status)
	}

	return r0
}

// mockWrapper_Status_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Status'
type mockWrapper_Status_Call struct {
	*mock.Call
}

// Status is a helper method to define mock.On call
func (_e *mockWrapper_Expecter) Status() *mockWrapper_Status_Call {
	return &mockWrapper_Status_Call{Call: _e.mock.On("Status")}
}

func (_c *mockWrapper_Status_Call) Run(run func()) *mockWrapper_Status_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockWrapper_Status_Call) Return(_a0 api.Status) *mockWrapper_Status_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockWrapper_Status_Call) RunAndReturn(run func() api.Status) *mockWrapper_Status_Call {
	_c.Call.Return(run)
	return _c
}

// newMockWrapper creates a new instance of mockWrapper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockWrapper(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockWrapper {
	mock := &mockWrapper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
