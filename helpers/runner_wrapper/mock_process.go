// Code generated by mockery v2.53.3. DO NOT EDIT.

package runner_wrapper

import (
	os "os"

	mock "github.com/stretchr/testify/mock"
)

// mockProcess is an autogenerated mock type for the process type
type mockProcess struct {
	mock.Mock
}

type mockProcess_Expecter struct {
	mock *mock.Mock
}

func (_m *mockProcess) EXPECT() *mockProcess_Expecter {
	return &mockProcess_Expecter{mock: &_m.Mock}
}

// Signal provides a mock function with given fields: sig
func (_m *mockProcess) Signal(sig os.Signal) error {
	ret := _m.Called(sig)

	if len(ret) == 0 {
		panic("no return value specified for Signal")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(os.Signal) error); ok {
		r0 = rf(sig)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockProcess_Signal_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Signal'
type mockProcess_Signal_Call struct {
	*mock.Call
}

// Signal is a helper method to define mock.On call
//   - sig os.Signal
func (_e *mockProcess_Expecter) Signal(sig interface{}) *mockProcess_Signal_Call {
	return &mockProcess_Signal_Call{Call: _e.mock.On("Signal", sig)}
}

func (_c *mockProcess_Signal_Call) Run(run func(sig os.Signal)) *mockProcess_Signal_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(os.Signal))
	})
	return _c
}

func (_c *mockProcess_Signal_Call) Return(_a0 error) *mockProcess_Signal_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockProcess_Signal_Call) RunAndReturn(run func(os.Signal) error) *mockProcess_Signal_Call {
	_c.Call.Return(run)
	return _c
}

// newMockProcess creates a new instance of mockProcess. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockProcess(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockProcess {
	mock := &mockProcess{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
