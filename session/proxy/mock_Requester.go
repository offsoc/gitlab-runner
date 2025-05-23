// Code generated by mockery v2.53.3. DO NOT EDIT.

package proxy

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// MockRequester is an autogenerated mock type for the Requester type
type MockRequester struct {
	mock.Mock
}

type MockRequester_Expecter struct {
	mock *mock.Mock
}

func (_m *MockRequester) EXPECT() *MockRequester_Expecter {
	return &MockRequester_Expecter{mock: &_m.Mock}
}

// ProxyRequest provides a mock function with given fields: w, r, requestedURI, port, settings
func (_m *MockRequester) ProxyRequest(w http.ResponseWriter, r *http.Request, requestedURI string, port string, settings *Settings) {
	_m.Called(w, r, requestedURI, port, settings)
}

// MockRequester_ProxyRequest_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProxyRequest'
type MockRequester_ProxyRequest_Call struct {
	*mock.Call
}

// ProxyRequest is a helper method to define mock.On call
//   - w http.ResponseWriter
//   - r *http.Request
//   - requestedURI string
//   - port string
//   - settings *Settings
func (_e *MockRequester_Expecter) ProxyRequest(w interface{}, r interface{}, requestedURI interface{}, port interface{}, settings interface{}) *MockRequester_ProxyRequest_Call {
	return &MockRequester_ProxyRequest_Call{Call: _e.mock.On("ProxyRequest", w, r, requestedURI, port, settings)}
}

func (_c *MockRequester_ProxyRequest_Call) Run(run func(w http.ResponseWriter, r *http.Request, requestedURI string, port string, settings *Settings)) *MockRequester_ProxyRequest_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(http.ResponseWriter), args[1].(*http.Request), args[2].(string), args[3].(string), args[4].(*Settings))
	})
	return _c
}

func (_c *MockRequester_ProxyRequest_Call) Return() *MockRequester_ProxyRequest_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockRequester_ProxyRequest_Call) RunAndReturn(run func(http.ResponseWriter, *http.Request, string, string, *Settings)) *MockRequester_ProxyRequest_Call {
	_c.Run(run)
	return _c
}

// NewMockRequester creates a new instance of MockRequester. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRequester(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRequester {
	mock := &MockRequester{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
