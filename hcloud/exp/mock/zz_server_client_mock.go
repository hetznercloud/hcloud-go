// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: IServerClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_server_client_mock.go -mock_names IServerClient=ServerClient github.com/hetznercloud/hcloud-go/v2/hcloud IServerClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// ServerClient is a mock of IServerClient interface.
type ServerClient struct {
	ctrl     *gomock.Controller
	recorder *ServerClientMockRecorder
}

// ServerClientMockRecorder is the mock recorder for ServerClient.
type ServerClientMockRecorder struct {
	mock *ServerClient
}

// NewServerClient creates a new mock instance.
func NewServerClient(ctrl *gomock.Controller) *ServerClient {
	mock := &ServerClient{ctrl: ctrl}
	mock.recorder = &ServerClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ServerClient) EXPECT() *ServerClientMockRecorder {
	return m.recorder
}

// AddToPlacementGroup mocks base method.
func (m *ServerClient) AddToPlacementGroup(arg0 context.Context, arg1 *hcloud.Server, arg2 *hcloud.PlacementGroup) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddToPlacementGroup", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddToPlacementGroup indicates an expected call of AddToPlacementGroup.
func (mr *ServerClientMockRecorder) AddToPlacementGroup(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddToPlacementGroup", reflect.TypeOf((*ServerClient)(nil).AddToPlacementGroup), arg0, arg1, arg2)
}

// All mocks base method.
func (m *ServerClient) All(arg0 context.Context) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *ServerClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*ServerClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *ServerClient) AllWithOpts(arg0 context.Context, arg1 hcloud.ServerListOpts) ([]*hcloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *ServerClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*ServerClient)(nil).AllWithOpts), arg0, arg1)
}

// AttachISO mocks base method.
func (m *ServerClient) AttachISO(arg0 context.Context, arg1 *hcloud.Server, arg2 *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachISO", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachISO indicates an expected call of AttachISO.
func (mr *ServerClientMockRecorder) AttachISO(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachISO", reflect.TypeOf((*ServerClient)(nil).AttachISO), arg0, arg1, arg2)
}

// AttachToNetwork mocks base method.
func (m *ServerClient) AttachToNetwork(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerAttachToNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachToNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachToNetwork indicates an expected call of AttachToNetwork.
func (mr *ServerClientMockRecorder) AttachToNetwork(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachToNetwork", reflect.TypeOf((*ServerClient)(nil).AttachToNetwork), arg0, arg1, arg2)
}

// ChangeAliasIPs mocks base method.
func (m *ServerClient) ChangeAliasIPs(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeAliasIPsOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAliasIPs", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeAliasIPs indicates an expected call of ChangeAliasIPs.
func (mr *ServerClientMockRecorder) ChangeAliasIPs(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAliasIPs", reflect.TypeOf((*ServerClient)(nil).ChangeAliasIPs), arg0, arg1, arg2)
}

// ChangeDNSPtr mocks base method.
func (m *ServerClient) ChangeDNSPtr(arg0 context.Context, arg1 *hcloud.Server, arg2 string, arg3 *string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeDNSPtr", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeDNSPtr indicates an expected call of ChangeDNSPtr.
func (mr *ServerClientMockRecorder) ChangeDNSPtr(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeDNSPtr", reflect.TypeOf((*ServerClient)(nil).ChangeDNSPtr), arg0, arg1, arg2, arg3)
}

// ChangeProtection mocks base method.
func (m *ServerClient) ChangeProtection(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeProtection", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeProtection indicates an expected call of ChangeProtection.
func (mr *ServerClientMockRecorder) ChangeProtection(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeProtection", reflect.TypeOf((*ServerClient)(nil).ChangeProtection), arg0, arg1, arg2)
}

// ChangeType mocks base method.
func (m *ServerClient) ChangeType(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerChangeTypeOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeType", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ChangeType indicates an expected call of ChangeType.
func (mr *ServerClientMockRecorder) ChangeType(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeType", reflect.TypeOf((*ServerClient)(nil).ChangeType), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *ServerClient) Create(arg0 context.Context, arg1 hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerCreateResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *ServerClientMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*ServerClient)(nil).Create), arg0, arg1)
}

// CreateImage mocks base method.
func (m *ServerClient) CreateImage(arg0 context.Context, arg1 *hcloud.Server, arg2 *hcloud.ServerCreateImageOpts) (hcloud.ServerCreateImageResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateImage", arg0, arg1, arg2)
	ret0, _ := ret[0].(hcloud.ServerCreateImageResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateImage indicates an expected call of CreateImage.
func (mr *ServerClientMockRecorder) CreateImage(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateImage", reflect.TypeOf((*ServerClient)(nil).CreateImage), arg0, arg1, arg2)
}

// Delete mocks base method.
func (m *ServerClient) Delete(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *ServerClientMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*ServerClient)(nil).Delete), arg0, arg1)
}

// DeleteWithResult mocks base method.
func (m *ServerClient) DeleteWithResult(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.ServerDeleteResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteWithResult", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.ServerDeleteResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteWithResult indicates an expected call of DeleteWithResult.
func (mr *ServerClientMockRecorder) DeleteWithResult(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteWithResult", reflect.TypeOf((*ServerClient)(nil).DeleteWithResult), arg0, arg1)
}

// DetachFromNetwork mocks base method.
func (m *ServerClient) DetachFromNetwork(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerDetachFromNetworkOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachFromNetwork", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachFromNetwork indicates an expected call of DetachFromNetwork.
func (mr *ServerClientMockRecorder) DetachFromNetwork(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachFromNetwork", reflect.TypeOf((*ServerClient)(nil).DetachFromNetwork), arg0, arg1, arg2)
}

// DetachISO mocks base method.
func (m *ServerClient) DetachISO(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachISO", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DetachISO indicates an expected call of DetachISO.
func (mr *ServerClientMockRecorder) DetachISO(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachISO", reflect.TypeOf((*ServerClient)(nil).DetachISO), arg0, arg1)
}

// DisableBackup mocks base method.
func (m *ServerClient) DisableBackup(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableBackup", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableBackup indicates an expected call of DisableBackup.
func (mr *ServerClientMockRecorder) DisableBackup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableBackup", reflect.TypeOf((*ServerClient)(nil).DisableBackup), arg0, arg1)
}

// DisableRescue mocks base method.
func (m *ServerClient) DisableRescue(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisableRescue", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DisableRescue indicates an expected call of DisableRescue.
func (mr *ServerClientMockRecorder) DisableRescue(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisableRescue", reflect.TypeOf((*ServerClient)(nil).DisableRescue), arg0, arg1)
}

// EnableBackup mocks base method.
func (m *ServerClient) EnableBackup(arg0 context.Context, arg1 *hcloud.Server, arg2 string) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableBackup", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableBackup indicates an expected call of EnableBackup.
func (mr *ServerClientMockRecorder) EnableBackup(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableBackup", reflect.TypeOf((*ServerClient)(nil).EnableBackup), arg0, arg1, arg2)
}

// EnableRescue mocks base method.
func (m *ServerClient) EnableRescue(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableRescue", arg0, arg1, arg2)
	ret0, _ := ret[0].(hcloud.ServerEnableRescueResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// EnableRescue indicates an expected call of EnableRescue.
func (mr *ServerClientMockRecorder) EnableRescue(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableRescue", reflect.TypeOf((*ServerClient)(nil).EnableRescue), arg0, arg1, arg2)
}

// Get mocks base method.
func (m *ServerClient) Get(arg0 context.Context, arg1 string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *ServerClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*ServerClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *ServerClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *ServerClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*ServerClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *ServerClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *ServerClientMockRecorder) GetByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*ServerClient)(nil).GetByName), arg0, arg1)
}

// GetMetrics mocks base method.
func (m *ServerClient) GetMetrics(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerGetMetricsOpts) (*hcloud.ServerMetrics, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.ServerMetrics)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *ServerClientMockRecorder) GetMetrics(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*ServerClient)(nil).GetMetrics), arg0, arg1, arg2)
}

// List mocks base method.
func (m *ServerClient) List(arg0 context.Context, arg1 hcloud.ServerListOpts) ([]*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *ServerClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*ServerClient)(nil).List), arg0, arg1)
}

// Poweroff mocks base method.
func (m *ServerClient) Poweroff(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweroff", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweroff indicates an expected call of Poweroff.
func (mr *ServerClientMockRecorder) Poweroff(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweroff", reflect.TypeOf((*ServerClient)(nil).Poweroff), arg0, arg1)
}

// Poweron mocks base method.
func (m *ServerClient) Poweron(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Poweron", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Poweron indicates an expected call of Poweron.
func (mr *ServerClientMockRecorder) Poweron(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Poweron", reflect.TypeOf((*ServerClient)(nil).Poweron), arg0, arg1)
}

// Reboot mocks base method.
func (m *ServerClient) Reboot(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reboot indicates an expected call of Reboot.
func (mr *ServerClientMockRecorder) Reboot(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*ServerClient)(nil).Reboot), arg0, arg1)
}

// Rebuild mocks base method.
func (m *ServerClient) Rebuild(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerRebuildOpts) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebuild", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Rebuild indicates an expected call of Rebuild.
func (mr *ServerClientMockRecorder) Rebuild(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebuild", reflect.TypeOf((*ServerClient)(nil).Rebuild), arg0, arg1, arg2)
}

// RebuildWithResult mocks base method.
func (m *ServerClient) RebuildWithResult(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerRebuildOpts) (hcloud.ServerRebuildResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RebuildWithResult", arg0, arg1, arg2)
	ret0, _ := ret[0].(hcloud.ServerRebuildResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RebuildWithResult indicates an expected call of RebuildWithResult.
func (mr *ServerClientMockRecorder) RebuildWithResult(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RebuildWithResult", reflect.TypeOf((*ServerClient)(nil).RebuildWithResult), arg0, arg1, arg2)
}

// RemoveFromPlacementGroup mocks base method.
func (m *ServerClient) RemoveFromPlacementGroup(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromPlacementGroup", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RemoveFromPlacementGroup indicates an expected call of RemoveFromPlacementGroup.
func (mr *ServerClientMockRecorder) RemoveFromPlacementGroup(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromPlacementGroup", reflect.TypeOf((*ServerClient)(nil).RemoveFromPlacementGroup), arg0, arg1)
}

// RequestConsole mocks base method.
func (m *ServerClient) RequestConsole(arg0 context.Context, arg1 *hcloud.Server) (hcloud.ServerRequestConsoleResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestConsole", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerRequestConsoleResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RequestConsole indicates an expected call of RequestConsole.
func (mr *ServerClientMockRecorder) RequestConsole(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestConsole", reflect.TypeOf((*ServerClient)(nil).RequestConsole), arg0, arg1)
}

// Reset mocks base method.
func (m *ServerClient) Reset(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reset", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Reset indicates an expected call of Reset.
func (mr *ServerClientMockRecorder) Reset(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*ServerClient)(nil).Reset), arg0, arg1)
}

// ResetPassword mocks base method.
func (m *ServerClient) ResetPassword(arg0 context.Context, arg1 *hcloud.Server) (hcloud.ServerResetPasswordResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", arg0, arg1)
	ret0, _ := ret[0].(hcloud.ServerResetPasswordResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *ServerClientMockRecorder) ResetPassword(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*ServerClient)(nil).ResetPassword), arg0, arg1)
}

// Shutdown mocks base method.
func (m *ServerClient) Shutdown(arg0 context.Context, arg1 *hcloud.Server) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Shutdown indicates an expected call of Shutdown.
func (mr *ServerClientMockRecorder) Shutdown(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*ServerClient)(nil).Shutdown), arg0, arg1)
}

// Update mocks base method.
func (m *ServerClient) Update(arg0 context.Context, arg1 *hcloud.Server, arg2 hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Server)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *ServerClientMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*ServerClient)(nil).Update), arg0, arg1, arg2)
}
