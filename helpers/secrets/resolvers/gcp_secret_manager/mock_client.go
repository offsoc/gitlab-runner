// Code generated by mockery v2.43.0. DO NOT EDIT.

package gcp_secret_manager

import (
	context "context"

	common "gitlab.com/gitlab-org/gitlab-runner/common"

	mock "github.com/stretchr/testify/mock"
)

// mockClient is an autogenerated mock type for the client type
type mockClient struct {
	mock.Mock
}

// GetSecret provides a mock function with given fields: ctx, s
func (_m *mockClient) GetSecret(ctx context.Context, s *common.GCPSecretManagerSecret) (string, error) {
	ret := _m.Called(ctx, s)

	if len(ret) == 0 {
		panic("no return value specified for GetSecret")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *common.GCPSecretManagerSecret) (string, error)); ok {
		return rf(ctx, s)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *common.GCPSecretManagerSecret) string); ok {
		r0 = rf(ctx, s)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *common.GCPSecretManagerSecret) error); ok {
		r1 = rf(ctx, s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// newMockClient creates a new instance of mockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockClient {
	mock := &mockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
