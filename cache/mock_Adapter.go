// Code generated by mockery v2.43.0. DO NOT EDIT.

package cache

import (
	context "context"
	url "net/url"

	mock "github.com/stretchr/testify/mock"
)

// MockAdapter is an autogenerated mock type for the Adapter type
type MockAdapter struct {
	mock.Mock
}

// GetDownloadURL provides a mock function with given fields: _a0
func (_m *MockAdapter) GetDownloadURL(_a0 context.Context) PresignedURL {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetDownloadURL")
	}

	var r0 PresignedURL
	if rf, ok := ret.Get(0).(func(context.Context) PresignedURL); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(PresignedURL)
	}

	return r0
}

// GetGoCloudURL provides a mock function with given fields: _a0
func (_m *MockAdapter) GetGoCloudURL(_a0 context.Context) *url.URL {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetGoCloudURL")
	}

	var r0 *url.URL
	if rf, ok := ret.Get(0).(func(context.Context) *url.URL); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*url.URL)
		}
	}

	return r0
}

// GetUploadEnv provides a mock function with given fields: _a0
func (_m *MockAdapter) GetUploadEnv(_a0 context.Context) map[string]string {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetUploadEnv")
	}

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context) map[string]string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	return r0
}

// GetUploadURL provides a mock function with given fields: _a0
func (_m *MockAdapter) GetUploadURL(_a0 context.Context) PresignedURL {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for GetUploadURL")
	}

	var r0 PresignedURL
	if rf, ok := ret.Get(0).(func(context.Context) PresignedURL); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(PresignedURL)
	}

	return r0
}

// NewMockAdapter creates a new instance of MockAdapter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAdapter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAdapter {
	mock := &MockAdapter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
