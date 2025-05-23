// Code generated by mockery v2.53.3. DO NOT EDIT.

package vault

import mock "github.com/stretchr/testify/mock"

// MockAuthMethod is an autogenerated mock type for the AuthMethod type
type MockAuthMethod struct {
	mock.Mock
}

type MockAuthMethod_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAuthMethod) EXPECT() *MockAuthMethod_Expecter {
	return &MockAuthMethod_Expecter{mock: &_m.Mock}
}

// Authenticate provides a mock function with given fields: client
func (_m *MockAuthMethod) Authenticate(client Client) error {
	ret := _m.Called(client)

	if len(ret) == 0 {
		panic("no return value specified for Authenticate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(Client) error); ok {
		r0 = rf(client)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockAuthMethod_Authenticate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Authenticate'
type MockAuthMethod_Authenticate_Call struct {
	*mock.Call
}

// Authenticate is a helper method to define mock.On call
//   - client Client
func (_e *MockAuthMethod_Expecter) Authenticate(client interface{}) *MockAuthMethod_Authenticate_Call {
	return &MockAuthMethod_Authenticate_Call{Call: _e.mock.On("Authenticate", client)}
}

func (_c *MockAuthMethod_Authenticate_Call) Run(run func(client Client)) *MockAuthMethod_Authenticate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(Client))
	})
	return _c
}

func (_c *MockAuthMethod_Authenticate_Call) Return(_a0 error) *MockAuthMethod_Authenticate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthMethod_Authenticate_Call) RunAndReturn(run func(Client) error) *MockAuthMethod_Authenticate_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with no fields
func (_m *MockAuthMethod) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockAuthMethod_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockAuthMethod_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockAuthMethod_Expecter) Name() *MockAuthMethod_Name_Call {
	return &MockAuthMethod_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockAuthMethod_Name_Call) Run(run func()) *MockAuthMethod_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAuthMethod_Name_Call) Return(_a0 string) *MockAuthMethod_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthMethod_Name_Call) RunAndReturn(run func() string) *MockAuthMethod_Name_Call {
	_c.Call.Return(run)
	return _c
}

// Token provides a mock function with no fields
func (_m *MockAuthMethod) Token() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Token")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockAuthMethod_Token_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Token'
type MockAuthMethod_Token_Call struct {
	*mock.Call
}

// Token is a helper method to define mock.On call
func (_e *MockAuthMethod_Expecter) Token() *MockAuthMethod_Token_Call {
	return &MockAuthMethod_Token_Call{Call: _e.mock.On("Token")}
}

func (_c *MockAuthMethod_Token_Call) Run(run func()) *MockAuthMethod_Token_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockAuthMethod_Token_Call) Return(_a0 string) *MockAuthMethod_Token_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockAuthMethod_Token_Call) RunAndReturn(run func() string) *MockAuthMethod_Token_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAuthMethod creates a new instance of MockAuthMethod. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuthMethod(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuthMethod {
	mock := &MockAuthMethod{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
