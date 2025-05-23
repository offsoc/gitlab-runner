// Code generated by mockery v2.53.3. DO NOT EDIT.

package wait

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockKillWaiter is an autogenerated mock type for the KillWaiter type
type MockKillWaiter struct {
	mock.Mock
}

type MockKillWaiter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockKillWaiter) EXPECT() *MockKillWaiter_Expecter {
	return &MockKillWaiter_Expecter{mock: &_m.Mock}
}

// StopKillWait provides a mock function with given fields: ctx, containerID, timeout, gracefulExitFunc
func (_m *MockKillWaiter) StopKillWait(ctx context.Context, containerID string, timeout *int, gracefulExitFunc GracefulExitFunc) error {
	ret := _m.Called(ctx, containerID, timeout, gracefulExitFunc)

	if len(ret) == 0 {
		panic("no return value specified for StopKillWait")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *int, GracefulExitFunc) error); ok {
		r0 = rf(ctx, containerID, timeout, gracefulExitFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKillWaiter_StopKillWait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopKillWait'
type MockKillWaiter_StopKillWait_Call struct {
	*mock.Call
}

// StopKillWait is a helper method to define mock.On call
//   - ctx context.Context
//   - containerID string
//   - timeout *int
//   - gracefulExitFunc GracefulExitFunc
func (_e *MockKillWaiter_Expecter) StopKillWait(ctx interface{}, containerID interface{}, timeout interface{}, gracefulExitFunc interface{}) *MockKillWaiter_StopKillWait_Call {
	return &MockKillWaiter_StopKillWait_Call{Call: _e.mock.On("StopKillWait", ctx, containerID, timeout, gracefulExitFunc)}
}

func (_c *MockKillWaiter_StopKillWait_Call) Run(run func(ctx context.Context, containerID string, timeout *int, gracefulExitFunc GracefulExitFunc)) *MockKillWaiter_StopKillWait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(*int), args[3].(GracefulExitFunc))
	})
	return _c
}

func (_c *MockKillWaiter_StopKillWait_Call) Return(_a0 error) *MockKillWaiter_StopKillWait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKillWaiter_StopKillWait_Call) RunAndReturn(run func(context.Context, string, *int, GracefulExitFunc) error) *MockKillWaiter_StopKillWait_Call {
	_c.Call.Return(run)
	return _c
}

// Wait provides a mock function with given fields: ctx, containerID
func (_m *MockKillWaiter) Wait(ctx context.Context, containerID string) error {
	ret := _m.Called(ctx, containerID)

	if len(ret) == 0 {
		panic("no return value specified for Wait")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockKillWaiter_Wait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Wait'
type MockKillWaiter_Wait_Call struct {
	*mock.Call
}

// Wait is a helper method to define mock.On call
//   - ctx context.Context
//   - containerID string
func (_e *MockKillWaiter_Expecter) Wait(ctx interface{}, containerID interface{}) *MockKillWaiter_Wait_Call {
	return &MockKillWaiter_Wait_Call{Call: _e.mock.On("Wait", ctx, containerID)}
}

func (_c *MockKillWaiter_Wait_Call) Run(run func(ctx context.Context, containerID string)) *MockKillWaiter_Wait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockKillWaiter_Wait_Call) Return(_a0 error) *MockKillWaiter_Wait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockKillWaiter_Wait_Call) RunAndReturn(run func(context.Context, string) error) *MockKillWaiter_Wait_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockKillWaiter creates a new instance of MockKillWaiter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockKillWaiter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockKillWaiter {
	mock := &MockKillWaiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
