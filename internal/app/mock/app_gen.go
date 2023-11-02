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

// CreatePlanet mocks base method.
func (m *MockApp) CreatePlanet(ctx context.Context, in input.CreateTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTicket", ctx, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlanet indicates an expected call of CreatePlanet.
func (mr *MockAppMockRecorder) CreatePlanet(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTicket", reflect.TypeOf((*MockApp)(nil).CreatePlanet), ctx, in)
}

// DeletePlanet mocks base method.
func (m *MockApp) DeletePlanet(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlanet", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlanet indicates an expected call of DeletePlanet.
func (mr *MockAppMockRecorder) DeletePlanet(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlanet", reflect.TypeOf((*MockApp)(nil).DeletePlanet), ctx, id)
}

// GetPlanet mocks base method.
func (m *MockApp) GetPlanet(ctx context.Context, id int64) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanet", ctx, id)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanet indicates an expected call of GetPlanet.
func (mr *MockAppMockRecorder) GetPlanet(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanet", reflect.TypeOf((*MockApp)(nil).GetPlanet), ctx, id)
}

// GetPlanetByCode mocks base method.
func (m *MockApp) GetPlanetByCode(ctx context.Context, code string) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanetByCode", ctx, code)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanetByCode indicates an expected call of GetPlanetByCode.
func (mr *MockAppMockRecorder) GetPlanetByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanetByCode", reflect.TypeOf((*MockApp)(nil).GetPlanetByCode), ctx, code)
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

// QueryPlanets mocks base method.
func (m *MockApp) QueryPlanets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryPlanets", ctx, query, page, perPage)
	ret0, _ := ret[0].(*pagination.Pages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryPlanets indicates an expected call of QueryPlanets.
func (mr *MockAppMockRecorder) QueryPlanets(ctx, query, page, perPage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryPlanets", reflect.TypeOf((*MockApp)(nil).QueryPlanets), ctx, query, page, perPage)
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

// UpdatePlanet mocks base method.
func (m *MockApp) UpdatePlanet(ctx context.Context, id int64, in input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicket", ctx, id, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePlanet indicates an expected call of UpdatePlanet.
func (mr *MockAppMockRecorder) UpdatePlanet(ctx, id, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicket", reflect.TypeOf((*MockApp)(nil).UpdatePlanet), ctx, id, in)
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

// MockPlanetService is a mock of PlanetService interface.
type MockPlanetService struct {
	ctrl     *gomock.Controller
	recorder *MockPlanetServiceMockRecorder
}

// MockPlanetServiceMockRecorder is the mock recorder for MockPlanetService.
type MockPlanetServiceMockRecorder struct {
	mock *MockPlanetService
}

// NewMockPlanetService creates a new mock instance.
func NewMockPlanetService(ctrl *gomock.Controller) *MockPlanetService {
	mock := &MockPlanetService{ctrl: ctrl}
	mock.recorder = &MockPlanetServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlanetService) EXPECT() *MockPlanetServiceMockRecorder {
	return m.recorder
}

// CreatePlanet mocks base method.
func (m *MockPlanetService) CreatePlanet(ctx context.Context, in input.CreateTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTicket", ctx, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlanet indicates an expected call of CreatePlanet.
func (mr *MockPlanetServiceMockRecorder) CreatePlanet(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTicket", reflect.TypeOf((*MockPlanetService)(nil).CreatePlanet), ctx, in)
}

// DeletePlanet mocks base method.
func (m *MockPlanetService) DeletePlanet(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlanet", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePlanet indicates an expected call of DeletePlanet.
func (mr *MockPlanetServiceMockRecorder) DeletePlanet(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlanet", reflect.TypeOf((*MockPlanetService)(nil).DeletePlanet), ctx, id)
}

// GetPlanet mocks base method.
func (m *MockPlanetService) GetPlanet(ctx context.Context, id int64) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanet", ctx, id)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanet indicates an expected call of GetPlanet.
func (mr *MockPlanetServiceMockRecorder) GetPlanet(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanet", reflect.TypeOf((*MockPlanetService)(nil).GetPlanet), ctx, id)
}

// GetPlanetByCode mocks base method.
func (m *MockPlanetService) GetPlanetByCode(ctx context.Context, code string) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanetByCode", ctx, code)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanetByCode indicates an expected call of GetPlanetByCode.
func (mr *MockPlanetServiceMockRecorder) GetPlanetByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanetByCode", reflect.TypeOf((*MockPlanetService)(nil).GetPlanetByCode), ctx, code)
}

// QueryPlanets mocks base method.
func (m *MockPlanetService) QueryPlanets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryPlanets", ctx, query, page, perPage)
	ret0, _ := ret[0].(*pagination.Pages)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryPlanets indicates an expected call of QueryPlanets.
func (mr *MockPlanetServiceMockRecorder) QueryPlanets(ctx, query, page, perPage interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryPlanets", reflect.TypeOf((*MockPlanetService)(nil).QueryPlanets), ctx, query, page, perPage)
}

// UpdatePlanet mocks base method.
func (m *MockPlanetService) UpdatePlanet(ctx context.Context, id int64, in input.PatchTicket) (*dto.Ticket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PatchTicket", ctx, id, in)
	ret0, _ := ret[0].(*dto.Ticket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePlanet indicates an expected call of UpdatePlanet.
func (mr *MockPlanetServiceMockRecorder) UpdatePlanet(ctx, id, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PatchTicket", reflect.TypeOf((*MockPlanetService)(nil).UpdatePlanet), ctx, id, in)
}
