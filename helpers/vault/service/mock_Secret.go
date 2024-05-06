// Code generated by mockery v2.43.0. DO NOT EDIT.

package service

import mock "github.com/stretchr/testify/mock"

// MockSecret is an autogenerated mock type for the Secret type
type MockSecret struct {
	mock.Mock
}

// SecretField provides a mock function with given fields:
func (_m *MockSecret) SecretField() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SecretField")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SecretPath provides a mock function with given fields:
func (_m *MockSecret) SecretPath() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for SecretPath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// NewMockSecret creates a new instance of MockSecret. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSecret(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSecret {
	mock := &MockSecret{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
