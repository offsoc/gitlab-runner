// Code generated by mockery v2.43.0. DO NOT EDIT.

package kubernetes

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// mockBackoffCalculator is an autogenerated mock type for the backoffCalculator type
type mockBackoffCalculator struct {
	mock.Mock
}

// ForAttempt provides a mock function with given fields: attempt
func (_m *mockBackoffCalculator) ForAttempt(attempt float64) time.Duration {
	ret := _m.Called(attempt)

	if len(ret) == 0 {
		panic("no return value specified for ForAttempt")
	}

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(float64) time.Duration); ok {
		r0 = rf(attempt)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// newMockBackoffCalculator creates a new instance of mockBackoffCalculator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockBackoffCalculator(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockBackoffCalculator {
	mock := &mockBackoffCalculator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
