// Code generated by mockery v2.53.3. DO NOT EDIT.

package api

import mock "github.com/stretchr/testify/mock"

// MockInitGracefulShutdownRequest is an autogenerated mock type for the InitGracefulShutdownRequest type
type MockInitGracefulShutdownRequest struct {
	mock.Mock
}

type MockInitGracefulShutdownRequest_Expecter struct {
	mock *mock.Mock
}

func (_m *MockInitGracefulShutdownRequest) EXPECT() *MockInitGracefulShutdownRequest_Expecter {
	return &MockInitGracefulShutdownRequest_Expecter{mock: &_m.Mock}
}

// ShutdownCallbackDef provides a mock function with no fields
func (_m *MockInitGracefulShutdownRequest) ShutdownCallbackDef() ShutdownCallbackDef {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ShutdownCallbackDef")
	}

	var r0 ShutdownCallbackDef
	if rf, ok := ret.Get(0).(func() ShutdownCallbackDef); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ShutdownCallbackDef)
		}
	}

	return r0
}

// MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ShutdownCallbackDef'
type MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call struct {
	*mock.Call
}

// ShutdownCallbackDef is a helper method to define mock.On call
func (_e *MockInitGracefulShutdownRequest_Expecter) ShutdownCallbackDef() *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call {
	return &MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call{Call: _e.mock.On("ShutdownCallbackDef")}
}

func (_c *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call) Run(run func()) *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call) Return(_a0 ShutdownCallbackDef) *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call) RunAndReturn(run func() ShutdownCallbackDef) *MockInitGracefulShutdownRequest_ShutdownCallbackDef_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockInitGracefulShutdownRequest creates a new instance of MockInitGracefulShutdownRequest. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInitGracefulShutdownRequest(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInitGracefulShutdownRequest {
	mock := &MockInitGracefulShutdownRequest{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
