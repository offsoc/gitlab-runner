// Code generated by mockery v2.53.3. DO NOT EDIT.

package referees

import (
	bytes "bytes"
	context "context"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockReferee is an autogenerated mock type for the Referee type
type MockReferee struct {
	mock.Mock
}

type MockReferee_Expecter struct {
	mock *mock.Mock
}

func (_m *MockReferee) EXPECT() *MockReferee_Expecter {
	return &MockReferee_Expecter{mock: &_m.Mock}
}

// ArtifactBaseName provides a mock function with no fields
func (_m *MockReferee) ArtifactBaseName() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArtifactBaseName")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockReferee_ArtifactBaseName_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArtifactBaseName'
type MockReferee_ArtifactBaseName_Call struct {
	*mock.Call
}

// ArtifactBaseName is a helper method to define mock.On call
func (_e *MockReferee_Expecter) ArtifactBaseName() *MockReferee_ArtifactBaseName_Call {
	return &MockReferee_ArtifactBaseName_Call{Call: _e.mock.On("ArtifactBaseName")}
}

func (_c *MockReferee_ArtifactBaseName_Call) Run(run func()) *MockReferee_ArtifactBaseName_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockReferee_ArtifactBaseName_Call) Return(_a0 string) *MockReferee_ArtifactBaseName_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockReferee_ArtifactBaseName_Call) RunAndReturn(run func() string) *MockReferee_ArtifactBaseName_Call {
	_c.Call.Return(run)
	return _c
}

// ArtifactFormat provides a mock function with no fields
func (_m *MockReferee) ArtifactFormat() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArtifactFormat")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockReferee_ArtifactFormat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArtifactFormat'
type MockReferee_ArtifactFormat_Call struct {
	*mock.Call
}

// ArtifactFormat is a helper method to define mock.On call
func (_e *MockReferee_Expecter) ArtifactFormat() *MockReferee_ArtifactFormat_Call {
	return &MockReferee_ArtifactFormat_Call{Call: _e.mock.On("ArtifactFormat")}
}

func (_c *MockReferee_ArtifactFormat_Call) Run(run func()) *MockReferee_ArtifactFormat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockReferee_ArtifactFormat_Call) Return(_a0 string) *MockReferee_ArtifactFormat_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockReferee_ArtifactFormat_Call) RunAndReturn(run func() string) *MockReferee_ArtifactFormat_Call {
	_c.Call.Return(run)
	return _c
}

// ArtifactType provides a mock function with no fields
func (_m *MockReferee) ArtifactType() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ArtifactType")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockReferee_ArtifactType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ArtifactType'
type MockReferee_ArtifactType_Call struct {
	*mock.Call
}

// ArtifactType is a helper method to define mock.On call
func (_e *MockReferee_Expecter) ArtifactType() *MockReferee_ArtifactType_Call {
	return &MockReferee_ArtifactType_Call{Call: _e.mock.On("ArtifactType")}
}

func (_c *MockReferee_ArtifactType_Call) Run(run func()) *MockReferee_ArtifactType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockReferee_ArtifactType_Call) Return(_a0 string) *MockReferee_ArtifactType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockReferee_ArtifactType_Call) RunAndReturn(run func() string) *MockReferee_ArtifactType_Call {
	_c.Call.Return(run)
	return _c
}

// Execute provides a mock function with given fields: ctx, startTime, endTime
func (_m *MockReferee) Execute(ctx context.Context, startTime time.Time, endTime time.Time) (*bytes.Reader, error) {
	ret := _m.Called(ctx, startTime, endTime)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *bytes.Reader
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) (*bytes.Reader, error)); ok {
		return rf(ctx, startTime, endTime)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, time.Time) *bytes.Reader); ok {
		r0 = rf(ctx, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bytes.Reader)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, time.Time) error); ok {
		r1 = rf(ctx, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockReferee_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockReferee_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - ctx context.Context
//   - startTime time.Time
//   - endTime time.Time
func (_e *MockReferee_Expecter) Execute(ctx interface{}, startTime interface{}, endTime interface{}) *MockReferee_Execute_Call {
	return &MockReferee_Execute_Call{Call: _e.mock.On("Execute", ctx, startTime, endTime)}
}

func (_c *MockReferee_Execute_Call) Run(run func(ctx context.Context, startTime time.Time, endTime time.Time)) *MockReferee_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Time), args[2].(time.Time))
	})
	return _c
}

func (_c *MockReferee_Execute_Call) Return(_a0 *bytes.Reader, _a1 error) *MockReferee_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockReferee_Execute_Call) RunAndReturn(run func(context.Context, time.Time, time.Time) (*bytes.Reader, error)) *MockReferee_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockReferee creates a new instance of MockReferee. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockReferee(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockReferee {
	mock := &MockReferee{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
