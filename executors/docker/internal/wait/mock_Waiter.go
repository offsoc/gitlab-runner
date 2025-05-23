// Code generated by mockery v2.53.3. DO NOT EDIT.

package wait

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockWaiter is an autogenerated mock type for the Waiter type
type MockWaiter struct {
	mock.Mock
}

type MockWaiter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWaiter) EXPECT() *MockWaiter_Expecter {
	return &MockWaiter_Expecter{mock: &_m.Mock}
}

// Wait provides a mock function with given fields: ctx, containerID
func (_m *MockWaiter) Wait(ctx context.Context, containerID string) error {
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

// MockWaiter_Wait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Wait'
type MockWaiter_Wait_Call struct {
	*mock.Call
}

// Wait is a helper method to define mock.On call
//   - ctx context.Context
//   - containerID string
func (_e *MockWaiter_Expecter) Wait(ctx interface{}, containerID interface{}) *MockWaiter_Wait_Call {
	return &MockWaiter_Wait_Call{Call: _e.mock.On("Wait", ctx, containerID)}
}

func (_c *MockWaiter_Wait_Call) Run(run func(ctx context.Context, containerID string)) *MockWaiter_Wait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockWaiter_Wait_Call) Return(_a0 error) *MockWaiter_Wait_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWaiter_Wait_Call) RunAndReturn(run func(context.Context, string) error) *MockWaiter_Wait_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWaiter creates a new instance of MockWaiter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWaiter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWaiter {
	mock := &MockWaiter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
