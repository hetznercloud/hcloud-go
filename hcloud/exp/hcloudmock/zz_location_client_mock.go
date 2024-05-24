// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: ILocationClient)
//
// Generated by this command:
//
//	mockgen -package hcloudmock -destination zz_location_client_mock.go -mock_names ILocationClient=LocationClient github.com/hetznercloud/hcloud-go/v2/hcloud ILocationClient
//

// Package hcloudmock is a generated GoMock package.
package hcloudmock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// LocationClient is a mock of ILocationClient interface.
type LocationClient struct {
	ctrl     *gomock.Controller
	recorder *LocationClientMockRecorder
}

// LocationClientMockRecorder is the mock recorder for LocationClient.
type LocationClientMockRecorder struct {
	mock *LocationClient
}

// NewLocationClient creates a new mock instance.
func NewLocationClient(ctrl *gomock.Controller) *LocationClient {
	mock := &LocationClient{ctrl: ctrl}
	mock.recorder = &LocationClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *LocationClient) EXPECT() *LocationClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *LocationClient) All(arg0 context.Context) ([]*hcloud.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *LocationClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*LocationClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *LocationClient) AllWithOpts(arg0 context.Context, arg1 hcloud.LocationListOpts) ([]*hcloud.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *LocationClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*LocationClient)(nil).AllWithOpts), arg0, arg1)
}

// Get mocks base method.
func (m *LocationClient) Get(arg0 context.Context, arg1 string) (*hcloud.Location, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Location)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *LocationClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*LocationClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *LocationClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Location, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Location)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *LocationClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*LocationClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *LocationClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Location, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Location)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *LocationClientMockRecorder) GetByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*LocationClient)(nil).GetByName), arg0, arg1)
}

// List mocks base method.
func (m *LocationClient) List(arg0 context.Context, arg1 hcloud.LocationListOpts) ([]*hcloud.Location, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Location)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *LocationClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*LocationClient)(nil).List), arg0, arg1)
}
