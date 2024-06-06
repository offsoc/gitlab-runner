// Code generated by mockery v2.43.0. DO NOT EDIT.

package log

import mock "github.com/stretchr/testify/mock"

// mockSystemLogger is an autogenerated mock type for the systemLogger type
type mockSystemLogger struct {
	mock.Mock
}

// Error provides a mock function with given fields: v
func (_m *mockSystemLogger) Error(v ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, v...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(v...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Errorf provides a mock function with given fields: format, a
func (_m *mockSystemLogger) Errorf(format string, a ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, a...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Errorf")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) error); ok {
		r0 = rf(format, a...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Info provides a mock function with given fields: v
func (_m *mockSystemLogger) Info(v ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, v...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Info")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(v...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Infof provides a mock function with given fields: format, a
func (_m *mockSystemLogger) Infof(format string, a ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, a...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Infof")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) error); ok {
		r0 = rf(format, a...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Warning provides a mock function with given fields: v
func (_m *mockSystemLogger) Warning(v ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, v...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Warning")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(v...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Warningf provides a mock function with given fields: format, a
func (_m *mockSystemLogger) Warningf(format string, a ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, a...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Warningf")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, ...interface{}) error); ok {
		r0 = rf(format, a...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockSystemLogger creates a new instance of mockSystemLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockSystemLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockSystemLogger {
	mock := &mockSystemLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
