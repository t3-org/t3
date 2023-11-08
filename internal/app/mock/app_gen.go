// Code generated by MockGen. DO NOT EDIT.
// Source: app_iface.go

// Package mockapp is a generated GoMock package.
package mockapp

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hexa "github.com/kamva/hexa"
	pagination "github.com/kamva/hexa/pagination"
	dto "space.org/space/internal/dto"
	input "space.org/space/internal/input"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// CreateTicket mocks base method.
func (m *MockApp) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTicket", ctx, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTicket indicates an expected call of CreateTicket.
func (mr *MockAppMockRecorder) CreateTicket(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTicket", reflect.TypeOf((*MockApp)(nil).CreateTicket), ctx, in)
}

// DeleteTicket mocks base method.
func (m *MockApp) DeleteTicket(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTicket", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTicket indicates an expected call of DeleteTicket.
func (mr *MockAppMockRecorder) DeleteTicket(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTicket", reflect.TypeOf((*MockApp)(nil).DeleteTicket), ctx, id)
}

// EditTicketUrlByKey mocks base method.
func (m *MockApp) EditTicketUrlByKey(ctx context.Context, key, val string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditTicketUrlByKey", ctx, key, val)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditTicketUrlByKey indicates an expected call of EditTicketUrlByKey.
func (mr *MockAppMockRecorder) EditTicketUrlByKey(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditTicketUrlByKey", reflect.TypeOf((*MockApp)(nil).EditTicketUrlByKey), ctx, key, val)
}

// GetTicket mocks base method.
func (m *MockApp) GetTicket(ctx context.Context, id int64) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicket", ctx, id)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTicket indicates an expected call of GetTicket.
func (mr *MockAppMockRecorder) GetTicket(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicket", reflect.TypeOf((*MockApp)(nil).GetTicket), ctx, id)
}

// GetTicketByKey mocks base method.
func (m *MockApp) GetTicketByKey(ctx context.Context, key, val string) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketByKey", ctx, key, val)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTicketByKey indicates an expected call of GetTicketByKey.
func (mr *MockAppMockRecorder) GetTicketByKey(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketByKey", reflect.TypeOf((*MockApp)(nil).GetTicketByKey), ctx, key, val)
}

// HealthIdentifier mocks base method.
func (m *MockApp) HealthIdentifier() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthIdentifier")
	ret0, _ := ret[0].(string)
	return ret0
}

// HealthIdentifier indicates an expected call of HealthIdentifier.
func (mr *MockAppMockRecorder) HealthIdentifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthIdentifier", reflect.TypeOf((*MockApp)(nil).HealthIdentifier))
}

// HealthStatus mocks base method.
func (m *MockApp) HealthStatus(ctx context.Context) hexa.HealthStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthStatus", ctx)
	ret0, _ := ret[0].(hexa.HealthStatus)
	return ret0
}

// HealthStatus indicates an expected call of HealthStatus.
func (mr *MockAppMockRecorder) HealthStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthStatus", reflect.TypeOf((*MockApp)(nil).HealthStatus), ctx)
}

// LivenessStatus mocks base method.
func (m *MockApp) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LivenessStatus", ctx)
	ret0, _ := ret[0].(hexa.LivenessStatus)
	return ret0
}

// LivenessStatus indicates an expected call of LivenessStatus.
func (mr *MockAppMockRecorder) LivenessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LivenessStatus", reflect.TypeOf((*MockApp)(nil).LivenessStatus), ctx)
}

// PatchTicket mocks base method.
func (m *MockApp) PatchTicket(ctx context.Context, id int64, in *input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicket", ctx, id, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTicket indicates an expected call of PatchTicket.
func (mr *MockAppMockRecorder) PatchTicket(ctx, id, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicket", reflect.TypeOf((*MockApp)(nil).PatchTicket), ctx, id, in)
}

// PatchTicketByKey mocks base method.
func (m *MockApp) PatchTicketByLabel(ctx context.Context, key, val string, in *input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicketByLabel", ctx, key, val, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTicketByKey indicates an expected call of PatchTicketByKey.
func (mr *MockAppMockRecorder) PatchTicketByKey(ctx, key, val, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicketByLabel", reflect.TypeOf((*MockApp)(nil).PatchTicketByLabel), ctx, key, val, in)
}

// QueryTickets mocks base method.
func (m *MockApp) QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTickets", ctx, query, page, perPage)
	ret0, _ := ret[0].(*pagination.Pages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTickets indicates an expected call of QueryTickets.
func (mr *MockAppMockRecorder) QueryTickets(ctx, query, page, perPage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTickets", reflect.TypeOf((*MockApp)(nil).QueryTickets), ctx, query, page, perPage)
}

// ReadinessStatus mocks base method.
func (m *MockApp) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadinessStatus", ctx)
	ret0, _ := ret[0].(hexa.ReadinessStatus)
	return ret0
}

// ReadinessStatus indicates an expected call of ReadinessStatus.
func (mr *MockAppMockRecorder) ReadinessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadinessStatus", reflect.TypeOf((*MockApp)(nil).ReadinessStatus), ctx)
}

// UpsertTickets mocks base method.
func (m *MockApp) UpsertTickets(ctx context.Context, in *input.BatchUpsertTickets) ([]*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertTickets", ctx, in)
	ret0, _ := ret[0].([]*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertTickets indicates an expected call of UpsertTickets.
func (mr *MockAppMockRecorder) UpsertTickets(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTickets", reflect.TypeOf((*MockApp)(nil).UpsertTickets), ctx, in)
}

// MockHealth is a mock of Health interface.
type MockHealth struct {
	ctrl     *gomock.Controller
	recorder *MockHealthMockRecorder
}

// MockHealthMockRecorder is the mock recorder for MockHealth.
type MockHealthMockRecorder struct {
	mock *MockHealth
}

// NewMockHealth creates a new mock instance.
func NewMockHealth(ctrl *gomock.Controller) *MockHealth {
	mock := &MockHealth{ctrl: ctrl}
	mock.recorder = &MockHealthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealth) EXPECT() *MockHealthMockRecorder {
	return m.recorder
}

// HealthIdentifier mocks base method.
func (m *MockHealth) HealthIdentifier() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthIdentifier")
	ret0, _ := ret[0].(string)
	return ret0
}

// HealthIdentifier indicates an expected call of HealthIdentifier.
func (mr *MockHealthMockRecorder) HealthIdentifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthIdentifier", reflect.TypeOf((*MockHealth)(nil).HealthIdentifier))
}

// HealthStatus mocks base method.
func (m *MockHealth) HealthStatus(ctx context.Context) hexa.HealthStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthStatus", ctx)
	ret0, _ := ret[0].(hexa.HealthStatus)
	return ret0
}

// HealthStatus indicates an expected call of HealthStatus.
func (mr *MockHealthMockRecorder) HealthStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthStatus", reflect.TypeOf((*MockHealth)(nil).HealthStatus), ctx)
}

// LivenessStatus mocks base method.
func (m *MockHealth) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LivenessStatus", ctx)
	ret0, _ := ret[0].(hexa.LivenessStatus)
	return ret0
}

// LivenessStatus indicates an expected call of LivenessStatus.
func (mr *MockHealthMockRecorder) LivenessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LivenessStatus", reflect.TypeOf((*MockHealth)(nil).LivenessStatus), ctx)
}

// ReadinessStatus mocks base method.
func (m *MockHealth) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadinessStatus", ctx)
	ret0, _ := ret[0].(hexa.ReadinessStatus)
	return ret0
}

// ReadinessStatus indicates an expected call of ReadinessStatus.
func (mr *MockHealthMockRecorder) ReadinessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadinessStatus", reflect.TypeOf((*MockHealth)(nil).ReadinessStatus), ctx)
}

// MockTicketService is a mock of TicketService interface.
type MockTicketService struct {
	ctrl     *gomock.Controller
	recorder *MockTicketServiceMockRecorder
}

// MockTicketServiceMockRecorder is the mock recorder for MockTicketService.
type MockTicketServiceMockRecorder struct {
	mock *MockTicketService
}

// NewMockTicketService creates a new mock instance.
func NewMockTicketService(ctrl *gomock.Controller) *MockTicketService {
	mock := &MockTicketService{ctrl: ctrl}
	mock.recorder = &MockTicketServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTicketService) EXPECT() *MockTicketServiceMockRecorder {
	return m.recorder
}

// CreateTicket mocks base method.
func (m *MockTicketService) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTicket", ctx, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTicket indicates an expected call of CreateTicket.
func (mr *MockTicketServiceMockRecorder) CreateTicket(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTicket", reflect.TypeOf((*MockTicketService)(nil).CreateTicket), ctx, in)
}

// DeleteTicket mocks base method.
func (m *MockTicketService) DeleteTicket(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTicket", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTicket indicates an expected call of DeleteTicket.
func (mr *MockTicketServiceMockRecorder) DeleteTicket(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTicket", reflect.TypeOf((*MockTicketService)(nil).DeleteTicket), ctx, id)
}

// EditTicketUrlByKey mocks base method.
func (m *MockTicketService) EditTicketUrlByKey(ctx context.Context, key, val string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditTicketUrlByKey", ctx, key, val)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditTicketUrlByKey indicates an expected call of EditTicketUrlByKey.
func (mr *MockTicketServiceMockRecorder) EditTicketUrlByKey(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditTicketUrlByKey", reflect.TypeOf((*MockTicketService)(nil).EditTicketUrlByKey), ctx, key, val)
}

// GetTicket mocks base method.
func (m *MockTicketService) GetTicket(ctx context.Context, id int64) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicket", ctx, id)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTicket indicates an expected call of GetTicket.
func (mr *MockTicketServiceMockRecorder) GetTicket(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicket", reflect.TypeOf((*MockTicketService)(nil).GetTicket), ctx, id)
}

// GetTicketByKey mocks base method.
func (m *MockTicketService) GetTicketByKey(ctx context.Context, key, val string) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTicketByKey", ctx, key, val)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTicketByKey indicates an expected call of GetTicketByKey.
func (mr *MockTicketServiceMockRecorder) GetTicketByKey(ctx, key, val interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTicketByKey", reflect.TypeOf((*MockTicketService)(nil).GetTicketByKey), ctx, key, val)
}

// PatchTicket mocks base method.
func (m *MockTicketService) PatchTicket(ctx context.Context, id int64, in *input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicket", ctx, id, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTicket indicates an expected call of PatchTicket.
func (mr *MockTicketServiceMockRecorder) PatchTicket(ctx, id, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicket", reflect.TypeOf((*MockTicketService)(nil).PatchTicket), ctx, id, in)
}

// PatchTicketByKey mocks base method.
func (m *MockTicketService) PatchTicketByLabel(ctx context.Context, key, val string, in *input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicketByLabel", ctx, key, val, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PatchTicketByKey indicates an expected call of PatchTicketByKey.
func (mr *MockTicketServiceMockRecorder) PatchTicketByKey(ctx, key, val, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicketByLabel", reflect.TypeOf((*MockTicketService)(nil).PatchTicketByLabel), ctx, key, val, in)
}

// QueryTickets mocks base method.
func (m *MockTicketService) QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTickets", ctx, query, page, perPage)
	ret0, _ := ret[0].(*pagination.Pages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTickets indicates an expected call of QueryTickets.
func (mr *MockTicketServiceMockRecorder) QueryTickets(ctx, query, page, perPage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTickets", reflect.TypeOf((*MockTicketService)(nil).QueryTickets), ctx, query, page, perPage)
}

// UpsertTickets mocks base method.
func (m *MockTicketService) UpsertTickets(ctx context.Context, in *input.BatchUpsertTickets) ([]*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertTickets", ctx, in)
	ret0, _ := ret[0].([]*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertTickets indicates an expected call of UpsertTickets.
func (mr *MockTicketServiceMockRecorder) UpsertTickets(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTickets", reflect.TypeOf((*MockTicketService)(nil).UpsertTickets), ctx, in)
}
