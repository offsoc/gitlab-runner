// Code generated by mockery v2.53.3. DO NOT EDIT.

package kubernetes

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
)

// mockPodStatusChecker is an autogenerated mock type for the podStatusChecker type
type mockPodStatusChecker struct {
	mock.Mock
}

type mockPodStatusChecker_Expecter struct {
	mock *mock.Mock
}

func (_m *mockPodStatusChecker) EXPECT() *mockPodStatusChecker_Expecter {
	return &mockPodStatusChecker_Expecter{mock: &_m.Mock}
}

// check provides a mock function with given fields: _a0, _a1
func (_m *mockPodStatusChecker) check(_a0 context.Context, _a1 *v1.Pod) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for check")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Pod) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockPodStatusChecker_check_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'check'
type mockPodStatusChecker_check_Call struct {
	*mock.Call
}

// check is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *v1.Pod
func (_e *mockPodStatusChecker_Expecter) check(_a0 interface{}, _a1 interface{}) *mockPodStatusChecker_check_Call {
	return &mockPodStatusChecker_check_Call{Call: _e.mock.On("check", _a0, _a1)}
}

func (_c *mockPodStatusChecker_check_Call) Run(run func(_a0 context.Context, _a1 *v1.Pod)) *mockPodStatusChecker_check_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*v1.Pod))
	})
	return _c
}

func (_c *mockPodStatusChecker_check_Call) Return(_a0 error) *mockPodStatusChecker_check_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockPodStatusChecker_check_Call) RunAndReturn(run func(context.Context, *v1.Pod) error) *mockPodStatusChecker_check_Call {
	_c.Call.Return(run)
	return _c
}

// newMockPodStatusChecker creates a new instance of mockPodStatusChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockPodStatusChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockPodStatusChecker {
	mock := &mockPodStatusChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
