// Code generated by mockery v2.43.0. DO NOT EDIT.

package docker

import (
	context "context"

	container "github.com/docker/docker/api/types/container"

	io "io"

	mock "github.com/stretchr/testify/mock"

	network "github.com/docker/docker/api/types/network"

	system "github.com/docker/docker/api/types/system"

	types "github.com/docker/docker/api/types"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"

	volume "github.com/docker/docker/api/types/volume"
)

// MockClient is an autogenerated mock type for the Client type
type MockClient struct {
	mock.Mock
}

// ClientVersion provides a mock function with given fields:
func (_m *MockClient) ClientVersion() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ClientVersion")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *MockClient) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerAttach provides a mock function with given fields: ctx, _a1, options
func (_m *MockClient) ContainerAttach(ctx context.Context, _a1 string, options container.AttachOptions) (types.HijackedResponse, error) {
	ret := _m.Called(ctx, _a1, options)

	if len(ret) == 0 {
		panic("no return value specified for ContainerAttach")
	}

	var r0 types.HijackedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.AttachOptions) (types.HijackedResponse, error)); ok {
		return rf(ctx, _a1, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, container.AttachOptions) types.HijackedResponse); ok {
		r0 = rf(ctx, _a1, options)
	} else {
		r0 = ret.Get(0).(types.HijackedResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, container.AttachOptions) error); ok {
		r1 = rf(ctx, _a1, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerCreate provides a mock function with given fields: ctx, config, hostConfig, networkingConfig, platform, containerName
func (_m *MockClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *v1.Platform, containerName string) (container.CreateResponse, error) {
	ret := _m.Called(ctx, config, hostConfig, networkingConfig, platform, containerName)

	if len(ret) == 0 {
		panic("no return value specified for ContainerCreate")
	}

	var r0 container.CreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) (container.CreateResponse, error)); ok {
		return rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) container.CreateResponse); ok {
		r0 = rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	} else {
		r0 = ret.Get(0).(container.CreateResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *container.Config, *container.HostConfig, *network.NetworkingConfig, *v1.Platform, string) error); ok {
		r1 = rf(ctx, config, hostConfig, networkingConfig, platform, containerName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerExecAttach provides a mock function with given fields: ctx, execID, config
func (_m *MockClient) ContainerExecAttach(ctx context.Context, execID string, config types.ExecStartCheck) (types.HijackedResponse, error) {
	ret := _m.Called(ctx, execID, config)

	if len(ret) == 0 {
		panic("no return value specified for ContainerExecAttach")
	}

	var r0 types.HijackedResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecStartCheck) (types.HijackedResponse, error)); ok {
		return rf(ctx, execID, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecStartCheck) types.HijackedResponse); ok {
		r0 = rf(ctx, execID, config)
	} else {
		r0 = ret.Get(0).(types.HijackedResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.ExecStartCheck) error); ok {
		r1 = rf(ctx, execID, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerExecCreate provides a mock function with given fields: ctx, _a1, config
func (_m *MockClient) ContainerExecCreate(ctx context.Context, _a1 string, config types.ExecConfig) (types.IDResponse, error) {
	ret := _m.Called(ctx, _a1, config)

	if len(ret) == 0 {
		panic("no return value specified for ContainerExecCreate")
	}

	var r0 types.IDResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecConfig) (types.IDResponse, error)); ok {
		return rf(ctx, _a1, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ExecConfig) types.IDResponse); ok {
		r0 = rf(ctx, _a1, config)
	} else {
		r0 = ret.Get(0).(types.IDResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.ExecConfig) error); ok {
		r1 = rf(ctx, _a1, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerInspect provides a mock function with given fields: ctx, containerID
func (_m *MockClient) ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	ret := _m.Called(ctx, containerID)

	if len(ret) == 0 {
		panic("no return value specified for ContainerInspect")
	}

	var r0 types.ContainerJSON
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.ContainerJSON, error)); ok {
		return rf(ctx, containerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.ContainerJSON); ok {
		r0 = rf(ctx, containerID)
	} else {
		r0 = ret.Get(0).(types.ContainerJSON)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, containerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerKill provides a mock function with given fields: ctx, containerID, signal
func (_m *MockClient) ContainerKill(ctx context.Context, containerID string, signal string) error {
	ret := _m.Called(ctx, containerID, signal)

	if len(ret) == 0 {
		panic("no return value specified for ContainerKill")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, containerID, signal)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerList provides a mock function with given fields: ctx, options
func (_m *MockClient) ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error) {
	ret := _m.Called(ctx, options)

	if len(ret) == 0 {
		panic("no return value specified for ContainerList")
	}

	var r0 []types.Container
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, container.ListOptions) ([]types.Container, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, container.ListOptions) []types.Container); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Container)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, container.ListOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerLogs provides a mock function with given fields: ctx, _a1, options
func (_m *MockClient) ContainerLogs(ctx context.Context, _a1 string, options container.LogsOptions) (io.ReadCloser, error) {
	ret := _m.Called(ctx, _a1, options)

	if len(ret) == 0 {
		panic("no return value specified for ContainerLogs")
	}

	var r0 io.ReadCloser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.LogsOptions) (io.ReadCloser, error)); ok {
		return rf(ctx, _a1, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, container.LogsOptions) io.ReadCloser); ok {
		r0 = rf(ctx, _a1, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, container.LogsOptions) error); ok {
		r1 = rf(ctx, _a1, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ContainerRemove provides a mock function with given fields: ctx, containerID, options
func (_m *MockClient) ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error {
	ret := _m.Called(ctx, containerID, options)

	if len(ret) == 0 {
		panic("no return value specified for ContainerRemove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.RemoveOptions) error); ok {
		r0 = rf(ctx, containerID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerStart provides a mock function with given fields: ctx, containerID, options
func (_m *MockClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	ret := _m.Called(ctx, containerID, options)

	if len(ret) == 0 {
		panic("no return value specified for ContainerStart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.StartOptions) error); ok {
		r0 = rf(ctx, containerID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerStop provides a mock function with given fields: ctx, containerID, opions
func (_m *MockClient) ContainerStop(ctx context.Context, containerID string, opions container.StopOptions) error {
	ret := _m.Called(ctx, containerID, opions)

	if len(ret) == 0 {
		panic("no return value specified for ContainerStop")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.StopOptions) error); ok {
		r0 = rf(ctx, containerID, opions)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ContainerWait provides a mock function with given fields: ctx, containerID, condition
func (_m *MockClient) ContainerWait(ctx context.Context, containerID string, condition container.WaitCondition) (<-chan container.WaitResponse, <-chan error) {
	ret := _m.Called(ctx, containerID, condition)

	if len(ret) == 0 {
		panic("no return value specified for ContainerWait")
	}

	var r0 <-chan container.WaitResponse
	var r1 <-chan error
	if rf, ok := ret.Get(0).(func(context.Context, string, container.WaitCondition) (<-chan container.WaitResponse, <-chan error)); ok {
		return rf(ctx, containerID, condition)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, container.WaitCondition) <-chan container.WaitResponse); ok {
		r0 = rf(ctx, containerID, condition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan container.WaitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, container.WaitCondition) <-chan error); ok {
		r1 = rf(ctx, containerID, condition)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}

// ImageImportBlocking provides a mock function with given fields: ctx, source, ref, options
func (_m *MockClient) ImageImportBlocking(ctx context.Context, source types.ImageImportSource, ref string, options types.ImageImportOptions) error {
	ret := _m.Called(ctx, source, ref, options)

	if len(ret) == 0 {
		panic("no return value specified for ImageImportBlocking")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, types.ImageImportSource, string, types.ImageImportOptions) error); ok {
		r0 = rf(ctx, source, ref, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ImageInspectWithRaw provides a mock function with given fields: ctx, imageID
func (_m *MockClient) ImageInspectWithRaw(ctx context.Context, imageID string) (types.ImageInspect, []byte, error) {
	ret := _m.Called(ctx, imageID)

	if len(ret) == 0 {
		panic("no return value specified for ImageInspectWithRaw")
	}

	var r0 types.ImageInspect
	var r1 []byte
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.ImageInspect, []byte, error)); ok {
		return rf(ctx, imageID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.ImageInspect); ok {
		r0 = rf(ctx, imageID)
	} else {
		r0 = ret.Get(0).(types.ImageInspect)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) []byte); ok {
		r1 = rf(ctx, imageID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, string) error); ok {
		r2 = rf(ctx, imageID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ImagePullBlocking provides a mock function with given fields: ctx, ref, options
func (_m *MockClient) ImagePullBlocking(ctx context.Context, ref string, options types.ImagePullOptions) error {
	ret := _m.Called(ctx, ref, options)

	if len(ret) == 0 {
		panic("no return value specified for ImagePullBlocking")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.ImagePullOptions) error); ok {
		r0 = rf(ctx, ref, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Info provides a mock function with given fields: ctx
func (_m *MockClient) Info(ctx context.Context) (system.Info, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Info")
	}

	var r0 system.Info
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (system.Info, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) system.Info); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(system.Info)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkCreate provides a mock function with given fields: ctx, networkName, options
func (_m *MockClient) NetworkCreate(ctx context.Context, networkName string, options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	ret := _m.Called(ctx, networkName, options)

	if len(ret) == 0 {
		panic("no return value specified for NetworkCreate")
	}

	var r0 types.NetworkCreateResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, types.NetworkCreate) (types.NetworkCreateResponse, error)); ok {
		return rf(ctx, networkName, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, types.NetworkCreate) types.NetworkCreateResponse); ok {
		r0 = rf(ctx, networkName, options)
	} else {
		r0 = ret.Get(0).(types.NetworkCreateResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, types.NetworkCreate) error); ok {
		r1 = rf(ctx, networkName, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkDisconnect provides a mock function with given fields: ctx, networkID, containerID, force
func (_m *MockClient) NetworkDisconnect(ctx context.Context, networkID string, containerID string, force bool) error {
	ret := _m.Called(ctx, networkID, containerID, force)

	if len(ret) == 0 {
		panic("no return value specified for NetworkDisconnect")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool) error); ok {
		r0 = rf(ctx, networkID, containerID, force)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NetworkInspect provides a mock function with given fields: ctx, networkID
func (_m *MockClient) NetworkInspect(ctx context.Context, networkID string) (types.NetworkResource, error) {
	ret := _m.Called(ctx, networkID)

	if len(ret) == 0 {
		panic("no return value specified for NetworkInspect")
	}

	var r0 types.NetworkResource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (types.NetworkResource, error)); ok {
		return rf(ctx, networkID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) types.NetworkResource); ok {
		r0 = rf(ctx, networkID)
	} else {
		r0 = ret.Get(0).(types.NetworkResource)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, networkID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkList provides a mock function with given fields: ctx, options
func (_m *MockClient) NetworkList(ctx context.Context, options types.NetworkListOptions) ([]types.NetworkResource, error) {
	ret := _m.Called(ctx, options)

	if len(ret) == 0 {
		panic("no return value specified for NetworkList")
	}

	var r0 []types.NetworkResource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, types.NetworkListOptions) ([]types.NetworkResource, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, types.NetworkListOptions) []types.NetworkResource); ok {
		r0 = rf(ctx, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.NetworkResource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, types.NetworkListOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NetworkRemove provides a mock function with given fields: ctx, networkID
func (_m *MockClient) NetworkRemove(ctx context.Context, networkID string) error {
	ret := _m.Called(ctx, networkID)

	if len(ret) == 0 {
		panic("no return value specified for NetworkRemove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, networkID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VolumeCreate provides a mock function with given fields: ctx, options
func (_m *MockClient) VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error) {
	ret := _m.Called(ctx, options)

	if len(ret) == 0 {
		panic("no return value specified for VolumeCreate")
	}

	var r0 volume.Volume
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, volume.CreateOptions) (volume.Volume, error)); ok {
		return rf(ctx, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, volume.CreateOptions) volume.Volume); ok {
		r0 = rf(ctx, options)
	} else {
		r0 = ret.Get(0).(volume.Volume)
	}

	if rf, ok := ret.Get(1).(func(context.Context, volume.CreateOptions) error); ok {
		r1 = rf(ctx, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VolumeInspect provides a mock function with given fields: ctx, volumeID
func (_m *MockClient) VolumeInspect(ctx context.Context, volumeID string) (volume.Volume, error) {
	ret := _m.Called(ctx, volumeID)

	if len(ret) == 0 {
		panic("no return value specified for VolumeInspect")
	}

	var r0 volume.Volume
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (volume.Volume, error)); ok {
		return rf(ctx, volumeID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) volume.Volume); ok {
		r0 = rf(ctx, volumeID)
	} else {
		r0 = ret.Get(0).(volume.Volume)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, volumeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VolumeRemove provides a mock function with given fields: ctx, volumeID, force
func (_m *MockClient) VolumeRemove(ctx context.Context, volumeID string, force bool) error {
	ret := _m.Called(ctx, volumeID, force)

	if len(ret) == 0 {
		panic("no return value specified for VolumeRemove")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) error); ok {
		r0 = rf(ctx, volumeID, force)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockClient creates a new instance of MockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockClient {
	mock := &MockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
