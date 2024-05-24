// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: IActionClient)
//
// Generated by this command:
//
//	mockgen -package hcloudmock -destination zz_action_client_mock.go -mock_names IActionClient=ActionClient github.com/hetznercloud/hcloud-go/v2/hcloud IActionClient
//

// Package hcloudmock is a generated GoMock package.
package hcloudmock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// ActionClient is a mock of IActionClient interface.
type ActionClient struct {
	ctrl     *gomock.Controller
	recorder *ActionClientMockRecorder
}

// ActionClientMockRecorder is the mock recorder for ActionClient.
type ActionClientMockRecorder struct {
	mock *ActionClient
}

// NewActionClient creates a new mock instance.
func NewActionClient(ctrl *gomock.Controller) *ActionClient {
	mock := &ActionClient{ctrl: ctrl}
	mock.recorder = &ActionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *ActionClient) EXPECT() *ActionClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *ActionClient) All(arg0 context.Context) ([]*hcloud.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *ActionClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*ActionClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *ActionClient) AllWithOpts(arg0 context.Context, arg1 hcloud.ActionListOpts) ([]*hcloud.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *ActionClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*ActionClient)(nil).AllWithOpts), arg0, arg1)
}

// GetByID mocks base method.
func (m *ActionClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *ActionClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*ActionClient)(nil).GetByID), arg0, arg1)
}

// List mocks base method.
func (m *ActionClient) List(arg0 context.Context, arg1 hcloud.ActionListOpts) ([]*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *ActionClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*ActionClient)(nil).List), arg0, arg1)
}

// WaitFor mocks base method.
func (m *ActionClient) WaitFor(arg0 context.Context, arg1 ...*hcloud.Action) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WaitFor", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitFor indicates an expected call of WaitFor.
func (mr *ActionClientMockRecorder) WaitFor(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitFor", reflect.TypeOf((*ActionClient)(nil).WaitFor), varargs...)
}

// WaitForFunc mocks base method.
func (m *ActionClient) WaitForFunc(arg0 context.Context, arg1 func(*hcloud.Action) error, arg2 ...*hcloud.Action) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WaitForFunc", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForFunc indicates an expected call of WaitForFunc.
func (mr *ActionClientMockRecorder) WaitForFunc(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForFunc", reflect.TypeOf((*ActionClient)(nil).WaitForFunc), varargs...)
}

// WatchOverallProgress mocks base method.
func (m *ActionClient) WatchOverallProgress(arg0 context.Context, arg1 []*hcloud.Action) (<-chan int, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchOverallProgress", arg0, arg1)
	ret0, _ := ret[0].(<-chan int)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// WatchOverallProgress indicates an expected call of WatchOverallProgress.
func (mr *ActionClientMockRecorder) WatchOverallProgress(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchOverallProgress", reflect.TypeOf((*ActionClient)(nil).WatchOverallProgress), arg0, arg1)
}

// WatchProgress mocks base method.
func (m *ActionClient) WatchProgress(arg0 context.Context, arg1 *hcloud.Action) (<-chan int, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchProgress", arg0, arg1)
	ret0, _ := ret[0].(<-chan int)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// WatchProgress indicates an expected call of WatchProgress.
func (mr *ActionClientMockRecorder) WatchProgress(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProgress", reflect.TypeOf((*ActionClient)(nil).WatchProgress), arg0, arg1)
}
