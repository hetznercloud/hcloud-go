// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hetznercloud/hcloud-go/v2/hcloud (interfaces: ICertificateClient)
//
// Generated by this command:
//
//	mockgen -package mock -destination zz_certificate_client_mock.go -mock_names ICertificateClient=MockCertificateClient github.com/hetznercloud/hcloud-go/v2/hcloud ICertificateClient
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	hcloud "github.com/hetznercloud/hcloud-go/v2/hcloud"
	gomock "go.uber.org/mock/gomock"
)

// MockCertificateClient is a mock of ICertificateClient interface.
type MockCertificateClient struct {
	ctrl     *gomock.Controller
	recorder *MockCertificateClientMockRecorder
}

// MockCertificateClientMockRecorder is the mock recorder for MockCertificateClient.
type MockCertificateClientMockRecorder struct {
	mock *MockCertificateClient
}

// NewMockCertificateClient creates a new mock instance.
func NewMockCertificateClient(ctrl *gomock.Controller) *MockCertificateClient {
	mock := &MockCertificateClient{ctrl: ctrl}
	mock.recorder = &MockCertificateClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCertificateClient) EXPECT() *MockCertificateClientMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockCertificateClient) All(arg0 context.Context) ([]*hcloud.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0)
	ret0, _ := ret[0].([]*hcloud.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockCertificateClientMockRecorder) All(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockCertificateClient)(nil).All), arg0)
}

// AllWithOpts mocks base method.
func (m *MockCertificateClient) AllWithOpts(arg0 context.Context, arg1 hcloud.CertificateListOpts) ([]*hcloud.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpts", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpts indicates an expected call of AllWithOpts.
func (mr *MockCertificateClientMockRecorder) AllWithOpts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpts", reflect.TypeOf((*MockCertificateClient)(nil).AllWithOpts), arg0, arg1)
}

// Create mocks base method.
func (m *MockCertificateClient) Create(arg0 context.Context, arg1 hcloud.CertificateCreateOpts) (*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockCertificateClientMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCertificateClient)(nil).Create), arg0, arg1)
}

// CreateCertificate mocks base method.
func (m *MockCertificateClient) CreateCertificate(arg0 context.Context, arg1 hcloud.CertificateCreateOpts) (hcloud.CertificateCreateResult, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCertificate", arg0, arg1)
	ret0, _ := ret[0].(hcloud.CertificateCreateResult)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateCertificate indicates an expected call of CreateCertificate.
func (mr *MockCertificateClientMockRecorder) CreateCertificate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCertificate", reflect.TypeOf((*MockCertificateClient)(nil).CreateCertificate), arg0, arg1)
}

// Delete mocks base method.
func (m *MockCertificateClient) Delete(arg0 context.Context, arg1 *hcloud.Certificate) (*hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockCertificateClientMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCertificateClient)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockCertificateClient) Get(arg0 context.Context, arg1 string) (*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockCertificateClientMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCertificateClient)(nil).Get), arg0, arg1)
}

// GetByID mocks base method.
func (m *MockCertificateClient) GetByID(arg0 context.Context, arg1 int64) (*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCertificateClientMockRecorder) GetByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCertificateClient)(nil).GetByID), arg0, arg1)
}

// GetByName mocks base method.
func (m *MockCertificateClient) GetByName(arg0 context.Context, arg1 string) (*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByName indicates an expected call of GetByName.
func (mr *MockCertificateClientMockRecorder) GetByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockCertificateClient)(nil).GetByName), arg0, arg1)
}

// List mocks base method.
func (m *MockCertificateClient) List(arg0 context.Context, arg1 hcloud.CertificateListOpts) ([]*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockCertificateClientMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCertificateClient)(nil).List), arg0, arg1)
}

// RetryIssuance mocks base method.
func (m *MockCertificateClient) RetryIssuance(arg0 context.Context, arg1 *hcloud.Certificate) (*hcloud.Action, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RetryIssuance", arg0, arg1)
	ret0, _ := ret[0].(*hcloud.Action)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RetryIssuance indicates an expected call of RetryIssuance.
func (mr *MockCertificateClientMockRecorder) RetryIssuance(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RetryIssuance", reflect.TypeOf((*MockCertificateClient)(nil).RetryIssuance), arg0, arg1)
}

// Update mocks base method.
func (m *MockCertificateClient) Update(arg0 context.Context, arg1 *hcloud.Certificate, arg2 hcloud.CertificateUpdateOpts) (*hcloud.Certificate, *hcloud.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(*hcloud.Certificate)
	ret1, _ := ret[1].(*hcloud.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockCertificateClientMockRecorder) Update(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCertificateClient)(nil).Update), arg0, arg1, arg2)
}
