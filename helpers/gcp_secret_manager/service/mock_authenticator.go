// Code generated by mockery v2.28.2. DO NOT EDIT.

package service

import (
	context "context"

	common "gitlab.com/gitlab-org/gitlab-runner/common"

	mock "github.com/stretchr/testify/mock"

	sts "google.golang.org/api/sts/v1"
)

// mockAuthenticator is an autogenerated mock type for the authenticator type
type mockAuthenticator struct {
	mock.Mock
}

// getToken provides a mock function with given fields: ctx, secret
func (_m *mockAuthenticator) getToken(ctx context.Context, secret *common.GCPSecretManagerSecret) (*sts.GoogleIdentityStsV1ExchangeTokenResponse, error) {
	ret := _m.Called(ctx, secret)

	var r0 *sts.GoogleIdentityStsV1ExchangeTokenResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *common.GCPSecretManagerSecret) (*sts.GoogleIdentityStsV1ExchangeTokenResponse, error)); ok {
		return rf(ctx, secret)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *common.GCPSecretManagerSecret) *sts.GoogleIdentityStsV1ExchangeTokenResponse); ok {
		r0 = rf(ctx, secret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sts.GoogleIdentityStsV1ExchangeTokenResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *common.GCPSecretManagerSecret) error); ok {
		r1 = rf(ctx, secret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTnewMockAuthenticator interface {
	mock.TestingT
	Cleanup(func())
}

// newMockAuthenticator creates a new instance of mockAuthenticator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockAuthenticator(t mockConstructorTestingTnewMockAuthenticator) *mockAuthenticator {
	mock := &mockAuthenticator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
