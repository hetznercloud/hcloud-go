// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: IImageClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_image_client_mock.go -mock_names IImageClient=MockImageClient github.com/hetznercloud/hcloud-go/v2/hcloud IImageClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// MockImageClient is a mock of IImageClient interface.
type MockImageClient struct {
	ctrl     *gomock.Controller
	recorder *MockImageClientMockRecorder
}

// MockImageClientMockRecorder is the mock recorder for MockImageClient.
type MockImageClientMockRecorder struct {
	mock *MockImageClient
}

// NewMockImageClient creates a new mock instance.
func NewMockImageClient(ctrl *gomock.Controller) *MockImageClient {
	mock := &MockImageClient{ctrl: ctrl}
	mock.recorder = &MockImageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageClient) EXPECT() *MockImageClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockImageClient) All(arg0 context.Context) ([]*hcloud.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockImageClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockImageClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockImageClient) AllWithOpts(arg0 context.Context, arg1 hcloud.ImageListOpts) ([]*hcloud.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockImageClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockImageClient)(nil).AllWithOpts), arg0, arg1)
}

// ChangeProtection mocks base method.
func (m *MockImageClient) ChangeProtection(arg0 context.Context, arg1 *hcloud.Image, arg2 hcloud.ImageChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeProtection", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeProtection indicates an expected call of ChangeProtection.
func (mr *MockImageClientMockRecorder) ChangeProtection(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeProtection", reflect.TypeOf((*MockImageClient)(nil).ChangeProtection), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *MockImageClient) Delete(arg0 context.Context, arg1 *hcloud.Image) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockImageClientMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockImageClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockImageClient) Get(arg0 context.Context, arg1 string) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockImageClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockImageClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockImageClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockImageClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockImageClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockImageClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockImageClientMockRecorder) GetByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockImageClient)(nil).GetByName), arg0, arg1)
}

// GetByNameAndArchitecture mocks base method.
func (m *MockImageClient) GetByNameAndArchitecture(arg0 context.Context, arg1 string, arg2 hcloud.Architecture) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByNameAndArchitecture", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByNameAndArchitecture indicates an expected call of GetByNameAndArchitecture.
func (mr *MockImageClientMockRecorder) GetByNameAndArchitecture(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByNameAndArchitecture", reflect.TypeOf((*MockImageClient)(nil).GetByNameAndArchitecture), arg0, arg1, arg2)
}

// GetForArchitecture mocks base method.
func (m *MockImageClient) GetForArchitecture(arg0 context.Context, arg1 string, arg2 hcloud.Architecture) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForArchitecture", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetForArchitecture indicates an expected call of GetForArchitecture.
func (mr *MockImageClientMockRecorder) GetForArchitecture(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForArchitecture", reflect.TypeOf((*MockImageClient)(nil).GetForArchitecture), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockImageClient) List(arg0 context.Context, arg1 hcloud.ImageListOpts) ([]*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockImageClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockImageClient)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockImageClient) Update(arg0 context.Context, arg1 *hcloud.Image, arg2 hcloud.ImageUpdateOpts) (*hcloud.Image, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Image)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockImageClientMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockImageClient)(nil).Update), arg0, arg1, arg2)
}
