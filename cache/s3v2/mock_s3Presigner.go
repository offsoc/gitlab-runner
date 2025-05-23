// Code generated by mockery v2.53.3. DO NOT EDIT.

package s3v2

import (
	context "context"

	cache "gitlab.com/gitlab-org/gitlab-runner/cache"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// mockS3Presigner is an autogenerated mock type for the s3Presigner type
type mockS3Presigner struct {
	mock.Mock
}

type mockS3Presigner_Expecter struct {
	mock *mock.Mock
}

func (_m *mockS3Presigner) EXPECT() *mockS3Presigner_Expecter {
	return &mockS3Presigner_Expecter{mock: &_m.Mock}
}

// FetchCredentialsForRole provides a mock function with given fields: ctx, roleARN, bucketName, objectName, upload, timeout
func (_m *mockS3Presigner) FetchCredentialsForRole(ctx context.Context, roleARN string, bucketName string, objectName string, upload bool, timeout time.Duration) (map[string]string, error) {
	ret := _m.Called(ctx, roleARN, bucketName, objectName, upload, timeout)

	if len(ret) == 0 {
		panic("no return value specified for FetchCredentialsForRole")
	}

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, bool, time.Duration) (map[string]string, error)); ok {
		return rf(ctx, roleARN, bucketName, objectName, upload, timeout)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, bool, time.Duration) map[string]string); ok {
		r0 = rf(ctx, roleARN, bucketName, objectName, upload, timeout)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, bool, time.Duration) error); ok {
		r1 = rf(ctx, roleARN, bucketName, objectName, upload, timeout)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockS3Presigner_FetchCredentialsForRole_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FetchCredentialsForRole'
type mockS3Presigner_FetchCredentialsForRole_Call struct {
	*mock.Call
}

// FetchCredentialsForRole is a helper method to define mock.On call
//   - ctx context.Context
//   - roleARN string
//   - bucketName string
//   - objectName string
//   - upload bool
//   - timeout time.Duration
func (_e *mockS3Presigner_Expecter) FetchCredentialsForRole(ctx interface{}, roleARN interface{}, bucketName interface{}, objectName interface{}, upload interface{}, timeout interface{}) *mockS3Presigner_FetchCredentialsForRole_Call {
	return &mockS3Presigner_FetchCredentialsForRole_Call{Call: _e.mock.On("FetchCredentialsForRole", ctx, roleARN, bucketName, objectName, upload, timeout)}
}

func (_c *mockS3Presigner_FetchCredentialsForRole_Call) Run(run func(ctx context.Context, roleARN string, bucketName string, objectName string, upload bool, timeout time.Duration)) *mockS3Presigner_FetchCredentialsForRole_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(bool), args[5].(time.Duration))
	})
	return _c
}

func (_c *mockS3Presigner_FetchCredentialsForRole_Call) Return(_a0 map[string]string, _a1 error) *mockS3Presigner_FetchCredentialsForRole_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockS3Presigner_FetchCredentialsForRole_Call) RunAndReturn(run func(context.Context, string, string, string, bool, time.Duration) (map[string]string, error)) *mockS3Presigner_FetchCredentialsForRole_Call {
	_c.Call.Return(run)
	return _c
}

// PresignURL provides a mock function with given fields: ctx, method, bucketName, objectName, expires
func (_m *mockS3Presigner) PresignURL(ctx context.Context, method string, bucketName string, objectName string, expires time.Duration) (cache.PresignedURL, error) {
	ret := _m.Called(ctx, method, bucketName, objectName, expires)

	if len(ret) == 0 {
		panic("no return value specified for PresignURL")
	}

	var r0 cache.PresignedURL
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, time.Duration) (cache.PresignedURL, error)); ok {
		return rf(ctx, method, bucketName, objectName, expires)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, time.Duration) cache.PresignedURL); ok {
		r0 = rf(ctx, method, bucketName, objectName, expires)
	} else {
		r0 = ret.Get(0).(cache.PresignedURL)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, time.Duration) error); ok {
		r1 = rf(ctx, method, bucketName, objectName, expires)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockS3Presigner_PresignURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PresignURL'
type mockS3Presigner_PresignURL_Call struct {
	*mock.Call
}

// PresignURL is a helper method to define mock.On call
//   - ctx context.Context
//   - method string
//   - bucketName string
//   - objectName string
//   - expires time.Duration
func (_e *mockS3Presigner_Expecter) PresignURL(ctx interface{}, method interface{}, bucketName interface{}, objectName interface{}, expires interface{}) *mockS3Presigner_PresignURL_Call {
	return &mockS3Presigner_PresignURL_Call{Call: _e.mock.On("PresignURL", ctx, method, bucketName, objectName, expires)}
}

func (_c *mockS3Presigner_PresignURL_Call) Run(run func(ctx context.Context, method string, bucketName string, objectName string, expires time.Duration)) *mockS3Presigner_PresignURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].(string), args[4].(time.Duration))
	})
	return _c
}

func (_c *mockS3Presigner_PresignURL_Call) Return(_a0 cache.PresignedURL, _a1 error) *mockS3Presigner_PresignURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockS3Presigner_PresignURL_Call) RunAndReturn(run func(context.Context, string, string, string, time.Duration) (cache.PresignedURL, error)) *mockS3Presigner_PresignURL_Call {
	_c.Call.Return(run)
	return _c
}

// ServerSideEncryptionType provides a mock function with no fields
func (_m *mockS3Presigner) ServerSideEncryptionType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ServerSideEncryptionType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// mockS3Presigner_ServerSideEncryptionType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ServerSideEncryptionType'
type mockS3Presigner_ServerSideEncryptionType_Call struct {
	*mock.Call
}

// ServerSideEncryptionType is a helper method to define mock.On call
func (_e *mockS3Presigner_Expecter) ServerSideEncryptionType() *mockS3Presigner_ServerSideEncryptionType_Call {
	return &mockS3Presigner_ServerSideEncryptionType_Call{Call: _e.mock.On("ServerSideEncryptionType")}
}

func (_c *mockS3Presigner_ServerSideEncryptionType_Call) Run(run func()) *mockS3Presigner_ServerSideEncryptionType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *mockS3Presigner_ServerSideEncryptionType_Call) Return(_a0 string) *mockS3Presigner_ServerSideEncryptionType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockS3Presigner_ServerSideEncryptionType_Call) RunAndReturn(run func() string) *mockS3Presigner_ServerSideEncryptionType_Call {
	_c.Call.Return(run)
	return _c
}

// newMockS3Presigner creates a new instance of mockS3Presigner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockS3Presigner(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockS3Presigner {
	mock := &mockS3Presigner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
