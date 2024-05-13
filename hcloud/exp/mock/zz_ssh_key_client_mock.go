// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: ISSHKeyClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_ssh_key_client_mock.go -mock_names ISSHKeyClient=SSHKeyClient github.com/hetznercloud/hcloud-go/v2/hcloud ISSHKeyClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// SSHKeyClient is a mock of ISSHKeyClient interface.
type SSHKeyClient struct {
	ctrl     *gomock.Controller
	recorder *SSHKeyClientMockRecorder
}

// SSHKeyClientMockRecorder is the mock recorder for SSHKeyClient.
type SSHKeyClientMockRecorder struct {
	mock *SSHKeyClient
}

// NewSSHKeyClient creates a new mock instance.
func NewSSHKeyClient(ctrl *gomock.Controller) *SSHKeyClient {
	mock := &SSHKeyClient{ctrl: ctrl}
	mock.recorder = &SSHKeyClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *SSHKeyClient) EXPECT() *SSHKeyClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *SSHKeyClient) All(arg0 context.Context) ([]*hcloud.SSHKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *SSHKeyClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*SSHKeyClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *SSHKeyClient) AllWithOpts(arg0 context.Context, arg1 hcloud.SSHKeyListOpts) ([]*hcloud.SSHKey, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *SSHKeyClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*SSHKeyClient)(nil).AllWithOpts), arg0, arg1)
}

// Create mocks base method.
func (m *SSHKeyClient) Create(arg0 context.Context, arg1 hcloud.SSHKeyCreateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *SSHKeyClientMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*SSHKeyClient)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *SSHKeyClient) Delete(arg0 context.Context, arg1 *hcloud.SSHKey) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *SSHKeyClientMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*SSHKeyClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *SSHKeyClient) Get(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *SSHKeyClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*SSHKeyClient)(nil).Get), arg0, arg1)
}

// GetByFingerprint mocks base method.
func (m *SSHKeyClient) GetByFingerprint(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByFingerprint", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByFingerprint indicates an expected call of GetByFingerprint.
func (mr *SSHKeyClientMockRecorder) GetByFingerprint(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByFingerprint", reflect.TypeOf((*SSHKeyClient)(nil).GetByFingerprint), arg0, arg1)
}

// GetByID mocks base method.
func (m *SSHKeyClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *SSHKeyClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*SSHKeyClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *SSHKeyClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *SSHKeyClientMockRecorder) GetByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*SSHKeyClient)(nil).GetByName), arg0, arg1)
}

// List mocks base method.
func (m *SSHKeyClient) List(arg0 context.Context, arg1 hcloud.SSHKeyListOpts) ([]*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *SSHKeyClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*SSHKeyClient)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *SSHKeyClient) Update(arg0 context.Context, arg1 *hcloud.SSHKey, arg2 hcloud.SSHKeyUpdateOpts) (*hcloud.SSHKey, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.SSHKey)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *SSHKeyClientMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*SSHKeyClient)(nil).Update), arg0, arg1, arg2)
}
