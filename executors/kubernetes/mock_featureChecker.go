// Code generated by mockery v2.43.0. DO NOT EDIT.

package kubernetes

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// mockFeatureChecker is an autogenerated mock type for the featureChecker type
type mockFeatureChecker struct {
	mock.Mock
}

// IsHostAliasSupported provides a mock function with given fields:
func (_m *mockFeatureChecker) IsHostAliasSupported() (bool, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsHostAliasSupported")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func() (bool, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsResourceVerbAllowed provides a mock function with given fields: _a0, _a1, _a2, _a3
func (_m *mockFeatureChecker) IsResourceVerbAllowed(_a0 context.Context, _a1 v1.GroupVersionResource, _a2 string, _a3 string) (bool, string, error) {
	ret := _m.Called(_a0, _a1, _a2, _a3)

	if len(ret) == 0 {
		panic("no return value specified for IsResourceVerbAllowed")
	}

	var r0 bool
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, v1.GroupVersionResource, string, string) (bool, string, error)); ok {
		return rf(_a0, _a1, _a2, _a3)
	}
	if rf, ok := ret.Get(0).(func(context.Context, v1.GroupVersionResource, string, string) bool); ok {
		r0 = rf(_a0, _a1, _a2, _a3)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, v1.GroupVersionResource, string, string) string); ok {
		r1 = rf(_a0, _a1, _a2, _a3)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(context.Context, v1.GroupVersionResource, string, string) error); ok {
		r2 = rf(_a0, _a1, _a2, _a3)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// newMockFeatureChecker creates a new instance of mockFeatureChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockFeatureChecker(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockFeatureChecker {
	mock := &mockFeatureChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
