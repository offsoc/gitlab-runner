// Code generated by mockery v2.53.3. DO NOT EDIT.

package common

import (
	context "context"
	io "io"

	mock "github.com/stretchr/testify/mock"

	time "time"
)

// MockNetwork is an autogenerated mock type for the Network type
type MockNetwork struct {
	mock.Mock
}

type MockNetwork_Expecter struct {
	mock *mock.Mock
}

func (_m *MockNetwork) EXPECT() *MockNetwork_Expecter {
	return &MockNetwork_Expecter{mock: &_m.Mock}
}

// DownloadArtifacts provides a mock function with given fields: config, artifactsFile, directDownload
func (_m *MockNetwork) DownloadArtifacts(config JobCredentials, artifactsFile io.WriteCloser, directDownload *bool) DownloadState {
	ret := _m.Called(config, artifactsFile, directDownload)

	if len(ret) == 0 {
		panic("no return value specified for DownloadArtifacts")
	}

	var r0 DownloadState
	if rf, ok := ret.Get(0).(func(JobCredentials, io.WriteCloser, *bool) DownloadState); ok {
		r0 = rf(config, artifactsFile, directDownload)
	} else {
		r0 = ret.Get(0).(DownloadState)
	}

	return r0
}

// MockNetwork_DownloadArtifacts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DownloadArtifacts'
type MockNetwork_DownloadArtifacts_Call struct {
	*mock.Call
}

// DownloadArtifacts is a helper method to define mock.On call
//   - config JobCredentials
//   - artifactsFile io.WriteCloser
//   - directDownload *bool
func (_e *MockNetwork_Expecter) DownloadArtifacts(config interface{}, artifactsFile interface{}, directDownload interface{}) *MockNetwork_DownloadArtifacts_Call {
	return &MockNetwork_DownloadArtifacts_Call{Call: _e.mock.On("DownloadArtifacts", config, artifactsFile, directDownload)}
}

func (_c *MockNetwork_DownloadArtifacts_Call) Run(run func(config JobCredentials, artifactsFile io.WriteCloser, directDownload *bool)) *MockNetwork_DownloadArtifacts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(JobCredentials), args[1].(io.WriteCloser), args[2].(*bool))
	})
	return _c
}

func (_c *MockNetwork_DownloadArtifacts_Call) Return(_a0 DownloadState) *MockNetwork_DownloadArtifacts_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_DownloadArtifacts_Call) RunAndReturn(run func(JobCredentials, io.WriteCloser, *bool) DownloadState) *MockNetwork_DownloadArtifacts_Call {
	_c.Call.Return(run)
	return _c
}

// PatchTrace provides a mock function with given fields: config, jobCredentials, content, startOffset, debugModeEnabled
func (_m *MockNetwork) PatchTrace(config RunnerConfig, jobCredentials *JobCredentials, content []byte, startOffset int, debugModeEnabled bool) PatchTraceResult {
	ret := _m.Called(config, jobCredentials, content, startOffset, debugModeEnabled)

	if len(ret) == 0 {
		panic("no return value specified for PatchTrace")
	}

	var r0 PatchTraceResult
	if rf, ok := ret.Get(0).(func(RunnerConfig, *JobCredentials, []byte, int, bool) PatchTraceResult); ok {
		r0 = rf(config, jobCredentials, content, startOffset, debugModeEnabled)
	} else {
		r0 = ret.Get(0).(PatchTraceResult)
	}

	return r0
}

// MockNetwork_PatchTrace_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PatchTrace'
type MockNetwork_PatchTrace_Call struct {
	*mock.Call
}

// PatchTrace is a helper method to define mock.On call
//   - config RunnerConfig
//   - jobCredentials *JobCredentials
//   - content []byte
//   - startOffset int
//   - debugModeEnabled bool
func (_e *MockNetwork_Expecter) PatchTrace(config interface{}, jobCredentials interface{}, content interface{}, startOffset interface{}, debugModeEnabled interface{}) *MockNetwork_PatchTrace_Call {
	return &MockNetwork_PatchTrace_Call{Call: _e.mock.On("PatchTrace", config, jobCredentials, content, startOffset, debugModeEnabled)}
}

func (_c *MockNetwork_PatchTrace_Call) Run(run func(config RunnerConfig, jobCredentials *JobCredentials, content []byte, startOffset int, debugModeEnabled bool)) *MockNetwork_PatchTrace_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerConfig), args[1].(*JobCredentials), args[2].([]byte), args[3].(int), args[4].(bool))
	})
	return _c
}

func (_c *MockNetwork_PatchTrace_Call) Return(_a0 PatchTraceResult) *MockNetwork_PatchTrace_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_PatchTrace_Call) RunAndReturn(run func(RunnerConfig, *JobCredentials, []byte, int, bool) PatchTraceResult) *MockNetwork_PatchTrace_Call {
	_c.Call.Return(run)
	return _c
}

// ProcessJob provides a mock function with given fields: config, buildCredentials
func (_m *MockNetwork) ProcessJob(config RunnerConfig, buildCredentials *JobCredentials) (JobTrace, error) {
	ret := _m.Called(config, buildCredentials)

	if len(ret) == 0 {
		panic("no return value specified for ProcessJob")
	}

	var r0 JobTrace
	var r1 error
	if rf, ok := ret.Get(0).(func(RunnerConfig, *JobCredentials) (JobTrace, error)); ok {
		return rf(config, buildCredentials)
	}
	if rf, ok := ret.Get(0).(func(RunnerConfig, *JobCredentials) JobTrace); ok {
		r0 = rf(config, buildCredentials)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(JobTrace)
		}
	}

	if rf, ok := ret.Get(1).(func(RunnerConfig, *JobCredentials) error); ok {
		r1 = rf(config, buildCredentials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockNetwork_ProcessJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ProcessJob'
type MockNetwork_ProcessJob_Call struct {
	*mock.Call
}

// ProcessJob is a helper method to define mock.On call
//   - config RunnerConfig
//   - buildCredentials *JobCredentials
func (_e *MockNetwork_Expecter) ProcessJob(config interface{}, buildCredentials interface{}) *MockNetwork_ProcessJob_Call {
	return &MockNetwork_ProcessJob_Call{Call: _e.mock.On("ProcessJob", config, buildCredentials)}
}

func (_c *MockNetwork_ProcessJob_Call) Run(run func(config RunnerConfig, buildCredentials *JobCredentials)) *MockNetwork_ProcessJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerConfig), args[1].(*JobCredentials))
	})
	return _c
}

func (_c *MockNetwork_ProcessJob_Call) Return(_a0 JobTrace, _a1 error) *MockNetwork_ProcessJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNetwork_ProcessJob_Call) RunAndReturn(run func(RunnerConfig, *JobCredentials) (JobTrace, error)) *MockNetwork_ProcessJob_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterRunner provides a mock function with given fields: config, parameters
func (_m *MockNetwork) RegisterRunner(config RunnerCredentials, parameters RegisterRunnerParameters) *RegisterRunnerResponse {
	ret := _m.Called(config, parameters)

	if len(ret) == 0 {
		panic("no return value specified for RegisterRunner")
	}

	var r0 *RegisterRunnerResponse
	if rf, ok := ret.Get(0).(func(RunnerCredentials, RegisterRunnerParameters) *RegisterRunnerResponse); ok {
		r0 = rf(config, parameters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RegisterRunnerResponse)
		}
	}

	return r0
}

// MockNetwork_RegisterRunner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterRunner'
type MockNetwork_RegisterRunner_Call struct {
	*mock.Call
}

// RegisterRunner is a helper method to define mock.On call
//   - config RunnerCredentials
//   - parameters RegisterRunnerParameters
func (_e *MockNetwork_Expecter) RegisterRunner(config interface{}, parameters interface{}) *MockNetwork_RegisterRunner_Call {
	return &MockNetwork_RegisterRunner_Call{Call: _e.mock.On("RegisterRunner", config, parameters)}
}

func (_c *MockNetwork_RegisterRunner_Call) Run(run func(config RunnerCredentials, parameters RegisterRunnerParameters)) *MockNetwork_RegisterRunner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials), args[1].(RegisterRunnerParameters))
	})
	return _c
}

func (_c *MockNetwork_RegisterRunner_Call) Return(_a0 *RegisterRunnerResponse) *MockNetwork_RegisterRunner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_RegisterRunner_Call) RunAndReturn(run func(RunnerCredentials, RegisterRunnerParameters) *RegisterRunnerResponse) *MockNetwork_RegisterRunner_Call {
	_c.Call.Return(run)
	return _c
}

// RequestJob provides a mock function with given fields: ctx, config, sessionInfo
func (_m *MockNetwork) RequestJob(ctx context.Context, config RunnerConfig, sessionInfo *SessionInfo) (*JobResponse, bool) {
	ret := _m.Called(ctx, config, sessionInfo)

	if len(ret) == 0 {
		panic("no return value specified for RequestJob")
	}

	var r0 *JobResponse
	var r1 bool
	if rf, ok := ret.Get(0).(func(context.Context, RunnerConfig, *SessionInfo) (*JobResponse, bool)); ok {
		return rf(ctx, config, sessionInfo)
	}
	if rf, ok := ret.Get(0).(func(context.Context, RunnerConfig, *SessionInfo) *JobResponse); ok {
		r0 = rf(ctx, config, sessionInfo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*JobResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, RunnerConfig, *SessionInfo) bool); ok {
		r1 = rf(ctx, config, sessionInfo)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// MockNetwork_RequestJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RequestJob'
type MockNetwork_RequestJob_Call struct {
	*mock.Call
}

// RequestJob is a helper method to define mock.On call
//   - ctx context.Context
//   - config RunnerConfig
//   - sessionInfo *SessionInfo
func (_e *MockNetwork_Expecter) RequestJob(ctx interface{}, config interface{}, sessionInfo interface{}) *MockNetwork_RequestJob_Call {
	return &MockNetwork_RequestJob_Call{Call: _e.mock.On("RequestJob", ctx, config, sessionInfo)}
}

func (_c *MockNetwork_RequestJob_Call) Run(run func(ctx context.Context, config RunnerConfig, sessionInfo *SessionInfo)) *MockNetwork_RequestJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(RunnerConfig), args[2].(*SessionInfo))
	})
	return _c
}

func (_c *MockNetwork_RequestJob_Call) Return(_a0 *JobResponse, _a1 bool) *MockNetwork_RequestJob_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNetwork_RequestJob_Call) RunAndReturn(run func(context.Context, RunnerConfig, *SessionInfo) (*JobResponse, bool)) *MockNetwork_RequestJob_Call {
	_c.Call.Return(run)
	return _c
}

// ResetToken provides a mock function with given fields: runner, systemID
func (_m *MockNetwork) ResetToken(runner RunnerCredentials, systemID string) *ResetTokenResponse {
	ret := _m.Called(runner, systemID)

	if len(ret) == 0 {
		panic("no return value specified for ResetToken")
	}

	var r0 *ResetTokenResponse
	if rf, ok := ret.Get(0).(func(RunnerCredentials, string) *ResetTokenResponse); ok {
		r0 = rf(runner, systemID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ResetTokenResponse)
		}
	}

	return r0
}

// MockNetwork_ResetToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetToken'
type MockNetwork_ResetToken_Call struct {
	*mock.Call
}

// ResetToken is a helper method to define mock.On call
//   - runner RunnerCredentials
//   - systemID string
func (_e *MockNetwork_Expecter) ResetToken(runner interface{}, systemID interface{}) *MockNetwork_ResetToken_Call {
	return &MockNetwork_ResetToken_Call{Call: _e.mock.On("ResetToken", runner, systemID)}
}

func (_c *MockNetwork_ResetToken_Call) Run(run func(runner RunnerCredentials, systemID string)) *MockNetwork_ResetToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials), args[1].(string))
	})
	return _c
}

func (_c *MockNetwork_ResetToken_Call) Return(_a0 *ResetTokenResponse) *MockNetwork_ResetToken_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_ResetToken_Call) RunAndReturn(run func(RunnerCredentials, string) *ResetTokenResponse) *MockNetwork_ResetToken_Call {
	_c.Call.Return(run)
	return _c
}

// ResetTokenWithPAT provides a mock function with given fields: runner, systemID, pat
func (_m *MockNetwork) ResetTokenWithPAT(runner RunnerCredentials, systemID string, pat string) *ResetTokenResponse {
	ret := _m.Called(runner, systemID, pat)

	if len(ret) == 0 {
		panic("no return value specified for ResetTokenWithPAT")
	}

	var r0 *ResetTokenResponse
	if rf, ok := ret.Get(0).(func(RunnerCredentials, string, string) *ResetTokenResponse); ok {
		r0 = rf(runner, systemID, pat)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ResetTokenResponse)
		}
	}

	return r0
}

// MockNetwork_ResetTokenWithPAT_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ResetTokenWithPAT'
type MockNetwork_ResetTokenWithPAT_Call struct {
	*mock.Call
}

// ResetTokenWithPAT is a helper method to define mock.On call
//   - runner RunnerCredentials
//   - systemID string
//   - pat string
func (_e *MockNetwork_Expecter) ResetTokenWithPAT(runner interface{}, systemID interface{}, pat interface{}) *MockNetwork_ResetTokenWithPAT_Call {
	return &MockNetwork_ResetTokenWithPAT_Call{Call: _e.mock.On("ResetTokenWithPAT", runner, systemID, pat)}
}

func (_c *MockNetwork_ResetTokenWithPAT_Call) Run(run func(runner RunnerCredentials, systemID string, pat string)) *MockNetwork_ResetTokenWithPAT_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockNetwork_ResetTokenWithPAT_Call) Return(_a0 *ResetTokenResponse) *MockNetwork_ResetTokenWithPAT_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_ResetTokenWithPAT_Call) RunAndReturn(run func(RunnerCredentials, string, string) *ResetTokenResponse) *MockNetwork_ResetTokenWithPAT_Call {
	_c.Call.Return(run)
	return _c
}

// SetConnectionMaxAge provides a mock function with given fields: _a0
func (_m *MockNetwork) SetConnectionMaxAge(_a0 time.Duration) {
	_m.Called(_a0)
}

// MockNetwork_SetConnectionMaxAge_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetConnectionMaxAge'
type MockNetwork_SetConnectionMaxAge_Call struct {
	*mock.Call
}

// SetConnectionMaxAge is a helper method to define mock.On call
//   - _a0 time.Duration
func (_e *MockNetwork_Expecter) SetConnectionMaxAge(_a0 interface{}) *MockNetwork_SetConnectionMaxAge_Call {
	return &MockNetwork_SetConnectionMaxAge_Call{Call: _e.mock.On("SetConnectionMaxAge", _a0)}
}

func (_c *MockNetwork_SetConnectionMaxAge_Call) Run(run func(_a0 time.Duration)) *MockNetwork_SetConnectionMaxAge_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *MockNetwork_SetConnectionMaxAge_Call) Return() *MockNetwork_SetConnectionMaxAge_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockNetwork_SetConnectionMaxAge_Call) RunAndReturn(run func(time.Duration)) *MockNetwork_SetConnectionMaxAge_Call {
	_c.Run(run)
	return _c
}

// UnregisterRunner provides a mock function with given fields: config
func (_m *MockNetwork) UnregisterRunner(config RunnerCredentials) bool {
	ret := _m.Called(config)

	if len(ret) == 0 {
		panic("no return value specified for UnregisterRunner")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(RunnerCredentials) bool); ok {
		r0 = rf(config)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockNetwork_UnregisterRunner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnregisterRunner'
type MockNetwork_UnregisterRunner_Call struct {
	*mock.Call
}

// UnregisterRunner is a helper method to define mock.On call
//   - config RunnerCredentials
func (_e *MockNetwork_Expecter) UnregisterRunner(config interface{}) *MockNetwork_UnregisterRunner_Call {
	return &MockNetwork_UnregisterRunner_Call{Call: _e.mock.On("UnregisterRunner", config)}
}

func (_c *MockNetwork_UnregisterRunner_Call) Run(run func(config RunnerCredentials)) *MockNetwork_UnregisterRunner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials))
	})
	return _c
}

func (_c *MockNetwork_UnregisterRunner_Call) Return(_a0 bool) *MockNetwork_UnregisterRunner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_UnregisterRunner_Call) RunAndReturn(run func(RunnerCredentials) bool) *MockNetwork_UnregisterRunner_Call {
	_c.Call.Return(run)
	return _c
}

// UnregisterRunnerManager provides a mock function with given fields: config, systemID
func (_m *MockNetwork) UnregisterRunnerManager(config RunnerCredentials, systemID string) bool {
	ret := _m.Called(config, systemID)

	if len(ret) == 0 {
		panic("no return value specified for UnregisterRunnerManager")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(RunnerCredentials, string) bool); ok {
		r0 = rf(config, systemID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// MockNetwork_UnregisterRunnerManager_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnregisterRunnerManager'
type MockNetwork_UnregisterRunnerManager_Call struct {
	*mock.Call
}

// UnregisterRunnerManager is a helper method to define mock.On call
//   - config RunnerCredentials
//   - systemID string
func (_e *MockNetwork_Expecter) UnregisterRunnerManager(config interface{}, systemID interface{}) *MockNetwork_UnregisterRunnerManager_Call {
	return &MockNetwork_UnregisterRunnerManager_Call{Call: _e.mock.On("UnregisterRunnerManager", config, systemID)}
}

func (_c *MockNetwork_UnregisterRunnerManager_Call) Run(run func(config RunnerCredentials, systemID string)) *MockNetwork_UnregisterRunnerManager_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials), args[1].(string))
	})
	return _c
}

func (_c *MockNetwork_UnregisterRunnerManager_Call) Return(_a0 bool) *MockNetwork_UnregisterRunnerManager_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_UnregisterRunnerManager_Call) RunAndReturn(run func(RunnerCredentials, string) bool) *MockNetwork_UnregisterRunnerManager_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateJob provides a mock function with given fields: config, jobCredentials, jobInfo
func (_m *MockNetwork) UpdateJob(config RunnerConfig, jobCredentials *JobCredentials, jobInfo UpdateJobInfo) UpdateJobResult {
	ret := _m.Called(config, jobCredentials, jobInfo)

	if len(ret) == 0 {
		panic("no return value specified for UpdateJob")
	}

	var r0 UpdateJobResult
	if rf, ok := ret.Get(0).(func(RunnerConfig, *JobCredentials, UpdateJobInfo) UpdateJobResult); ok {
		r0 = rf(config, jobCredentials, jobInfo)
	} else {
		r0 = ret.Get(0).(UpdateJobResult)
	}

	return r0
}

// MockNetwork_UpdateJob_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateJob'
type MockNetwork_UpdateJob_Call struct {
	*mock.Call
}

// UpdateJob is a helper method to define mock.On call
//   - config RunnerConfig
//   - jobCredentials *JobCredentials
//   - jobInfo UpdateJobInfo
func (_e *MockNetwork_Expecter) UpdateJob(config interface{}, jobCredentials interface{}, jobInfo interface{}) *MockNetwork_UpdateJob_Call {
	return &MockNetwork_UpdateJob_Call{Call: _e.mock.On("UpdateJob", config, jobCredentials, jobInfo)}
}

func (_c *MockNetwork_UpdateJob_Call) Run(run func(config RunnerConfig, jobCredentials *JobCredentials, jobInfo UpdateJobInfo)) *MockNetwork_UpdateJob_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerConfig), args[1].(*JobCredentials), args[2].(UpdateJobInfo))
	})
	return _c
}

func (_c *MockNetwork_UpdateJob_Call) Return(_a0 UpdateJobResult) *MockNetwork_UpdateJob_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_UpdateJob_Call) RunAndReturn(run func(RunnerConfig, *JobCredentials, UpdateJobInfo) UpdateJobResult) *MockNetwork_UpdateJob_Call {
	_c.Call.Return(run)
	return _c
}

// UploadRawArtifacts provides a mock function with given fields: config, bodyProvider, options
func (_m *MockNetwork) UploadRawArtifacts(config JobCredentials, bodyProvider ContentProvider, options ArtifactsOptions) (UploadState, string) {
	ret := _m.Called(config, bodyProvider, options)

	if len(ret) == 0 {
		panic("no return value specified for UploadRawArtifacts")
	}

	var r0 UploadState
	var r1 string
	if rf, ok := ret.Get(0).(func(JobCredentials, ContentProvider, ArtifactsOptions) (UploadState, string)); ok {
		return rf(config, bodyProvider, options)
	}
	if rf, ok := ret.Get(0).(func(JobCredentials, ContentProvider, ArtifactsOptions) UploadState); ok {
		r0 = rf(config, bodyProvider, options)
	} else {
		r0 = ret.Get(0).(UploadState)
	}

	if rf, ok := ret.Get(1).(func(JobCredentials, ContentProvider, ArtifactsOptions) string); ok {
		r1 = rf(config, bodyProvider, options)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// MockNetwork_UploadRawArtifacts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UploadRawArtifacts'
type MockNetwork_UploadRawArtifacts_Call struct {
	*mock.Call
}

// UploadRawArtifacts is a helper method to define mock.On call
//   - config JobCredentials
//   - bodyProvider ContentProvider
//   - options ArtifactsOptions
func (_e *MockNetwork_Expecter) UploadRawArtifacts(config interface{}, bodyProvider interface{}, options interface{}) *MockNetwork_UploadRawArtifacts_Call {
	return &MockNetwork_UploadRawArtifacts_Call{Call: _e.mock.On("UploadRawArtifacts", config, bodyProvider, options)}
}

func (_c *MockNetwork_UploadRawArtifacts_Call) Run(run func(config JobCredentials, bodyProvider ContentProvider, options ArtifactsOptions)) *MockNetwork_UploadRawArtifacts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(JobCredentials), args[1].(ContentProvider), args[2].(ArtifactsOptions))
	})
	return _c
}

func (_c *MockNetwork_UploadRawArtifacts_Call) Return(_a0 UploadState, _a1 string) *MockNetwork_UploadRawArtifacts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockNetwork_UploadRawArtifacts_Call) RunAndReturn(run func(JobCredentials, ContentProvider, ArtifactsOptions) (UploadState, string)) *MockNetwork_UploadRawArtifacts_Call {
	_c.Call.Return(run)
	return _c
}

// VerifyRunner provides a mock function with given fields: config, systemID
func (_m *MockNetwork) VerifyRunner(config RunnerCredentials, systemID string) *VerifyRunnerResponse {
	ret := _m.Called(config, systemID)

	if len(ret) == 0 {
		panic("no return value specified for VerifyRunner")
	}

	var r0 *VerifyRunnerResponse
	if rf, ok := ret.Get(0).(func(RunnerCredentials, string) *VerifyRunnerResponse); ok {
		r0 = rf(config, systemID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*VerifyRunnerResponse)
		}
	}

	return r0
}

// MockNetwork_VerifyRunner_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'VerifyRunner'
type MockNetwork_VerifyRunner_Call struct {
	*mock.Call
}

// VerifyRunner is a helper method to define mock.On call
//   - config RunnerCredentials
//   - systemID string
func (_e *MockNetwork_Expecter) VerifyRunner(config interface{}, systemID interface{}) *MockNetwork_VerifyRunner_Call {
	return &MockNetwork_VerifyRunner_Call{Call: _e.mock.On("VerifyRunner", config, systemID)}
}

func (_c *MockNetwork_VerifyRunner_Call) Run(run func(config RunnerCredentials, systemID string)) *MockNetwork_VerifyRunner_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(RunnerCredentials), args[1].(string))
	})
	return _c
}

func (_c *MockNetwork_VerifyRunner_Call) Return(_a0 *VerifyRunnerResponse) *MockNetwork_VerifyRunner_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockNetwork_VerifyRunner_Call) RunAndReturn(run func(RunnerCredentials, string) *VerifyRunnerResponse) *MockNetwork_VerifyRunner_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockNetwork creates a new instance of MockNetwork. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockNetwork(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockNetwork {
	mock := &MockNetwork{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
