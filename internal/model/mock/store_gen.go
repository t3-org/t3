// Code generated by MockGen. DO NOT EDIT.
// Source: store.go

// Package mockmodel is a generated GoMock package.
package mockmodel

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hexa "github.com/kamva/hexa"
	model "space.org/space/internal/model"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// DBLayer mocks base method.
func (m *MockStore) DBLayer() model.Store {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DBLayer")
	ret0, _ := ret[0].(model.Store)
	return ret0
}

// DBLayer indicates an expected call of DBLayer.
func (mr *MockStoreMockRecorder) DBLayer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DBLayer", reflect.TypeOf((*MockStore)(nil).DBLayer))
}

// HealthIdentifier mocks base method.
func (m *MockStore) HealthIdentifier() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthIdentifier")
	ret0, _ := ret[0].(string)
	return ret0
}

// HealthIdentifier indicates an expected call of HealthIdentifier.
func (mr *MockStoreMockRecorder) HealthIdentifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthIdentifier", reflect.TypeOf((*MockStore)(nil).HealthIdentifier))
}

// HealthStatus mocks base method.
func (m *MockStore) HealthStatus(ctx context.Context) hexa.HealthStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HealthStatus", ctx)
	ret0, _ := ret[0].(hexa.HealthStatus)
	return ret0
}

// HealthStatus indicates an expected call of HealthStatus.
func (mr *MockStoreMockRecorder) HealthStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HealthStatus", reflect.TypeOf((*MockStore)(nil).HealthStatus), ctx)
}

// LivenessStatus mocks base method.
func (m *MockStore) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LivenessStatus", ctx)
	ret0, _ := ret[0].(hexa.LivenessStatus)
	return ret0
}

// LivenessStatus indicates an expected call of LivenessStatus.
func (mr *MockStoreMockRecorder) LivenessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LivenessStatus", reflect.TypeOf((*MockStore)(nil).LivenessStatus), ctx)
}

// Planet mocks base method.
func (m *MockStore) Planet() model.PlanetStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Planet")
	ret0, _ := ret[0].(model.PlanetStore)
	return ret0
}

// Planet indicates an expected call of Planet.
func (mr *MockStoreMockRecorder) Planet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Planet", reflect.TypeOf((*MockStore)(nil).Planet))
}

// ReadinessStatus mocks base method.
func (m *MockStore) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadinessStatus", ctx)
	ret0, _ := ret[0].(hexa.ReadinessStatus)
	return ret0
}

// ReadinessStatus indicates an expected call of ReadinessStatus.
func (mr *MockStoreMockRecorder) ReadinessStatus(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadinessStatus", reflect.TypeOf((*MockStore)(nil).ReadinessStatus), ctx)
}

// System mocks base method.
func (m *MockStore) System() model.SystemStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "System")
	ret0, _ := ret[0].(model.SystemStore)
	return ret0
}

// System indicates an expected call of System.
func (mr *MockStoreMockRecorder) System() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "System", reflect.TypeOf((*MockStore)(nil).System))
}

// TruncateAllTables mocks base method.
func (m *MockStore) TruncateAllTables(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TruncateAllTables", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// TruncateAllTables indicates an expected call of TruncateAllTables.
func (mr *MockStoreMockRecorder) TruncateAllTables(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TruncateAllTables", reflect.TypeOf((*MockStore)(nil).TruncateAllTables), ctx)
}

// MockSystemStore is a mock of SystemStore interface.
type MockSystemStore struct {
	ctrl     *gomock.Controller
	recorder *MockSystemStoreMockRecorder
}

// MockSystemStoreMockRecorder is the mock recorder for MockSystemStore.
type MockSystemStoreMockRecorder struct {
	mock *MockSystemStore
}

// NewMockSystemStore creates a new mock instance.
func NewMockSystemStore(ctrl *gomock.Controller) *MockSystemStore {
	mock := &MockSystemStore{ctrl: ctrl}
	mock.recorder = &MockSystemStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystemStore) EXPECT() *MockSystemStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSystemStore) Delete(ctx context.Context, name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSystemStoreMockRecorder) Delete(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSystemStore)(nil).Delete), ctx, name)
}

// GetByName mocks base method.
func (m *MockSystemStore) GetByName(ctx context.Context, name string) (*model.System, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(*model.System)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockSystemStoreMockRecorder) GetByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockSystemStore)(nil).GetByName), ctx, name)
}

// Save mocks base method.
func (m *MockSystemStore) Save(ctx context.Context, system *model.System) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, system)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockSystemStoreMockRecorder) Save(ctx, system interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSystemStore)(nil).Save), ctx, system)
}

// MockPlanetStore is a mock of PlanetStore interface.
type MockPlanetStore struct {
	ctrl     *gomock.Controller
	recorder *MockPlanetStoreMockRecorder
}

// MockPlanetStoreMockRecorder is the mock recorder for MockPlanetStore.
type MockPlanetStoreMockRecorder struct {
	mock *MockPlanetStore
}

// NewMockPlanetStore creates a new mock instance.
func NewMockPlanetStore(ctrl *gomock.Controller) *MockPlanetStore {
	mock := &MockPlanetStore{ctrl: ctrl}
	mock.recorder = &MockPlanetStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlanetStore) EXPECT() *MockPlanetStoreMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockPlanetStore) Count(ctx context.Context, query string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx, query)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockPlanetStoreMockRecorder) Count(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockPlanetStore)(nil).Count), ctx, query)
}

// Create mocks base method.
func (m_2 *MockPlanetStore) Create(ctx context.Context, m *model.Planet) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Create", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockPlanetStoreMockRecorder) Create(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPlanetStore)(nil).Create), ctx, m)
}

// Delete mocks base method.
func (m_2 *MockPlanetStore) Delete(ctx context.Context, m *model.Planet) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Delete", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPlanetStoreMockRecorder) Delete(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPlanetStore)(nil).Delete), ctx, m)
}

// Get mocks base method.
func (m *MockPlanetStore) Get(ctx context.Context, id int64) (*model.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPlanetStoreMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPlanetStore)(nil).Get), ctx, id)
}

// GetByCode mocks base method.
func (m *MockPlanetStore) GetByCode(ctx context.Context, code string) (*model.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCode", ctx, code)
	ret0, _ := ret[0].(*model.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCode indicates an expected call of GetByCode.
func (mr *MockPlanetStoreMockRecorder) GetByCode(ctx, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCode", reflect.TypeOf((*MockPlanetStore)(nil).GetByCode), ctx, code)
}

// Query mocks base method.
func (m *MockPlanetStore) Query(ctx context.Context, query string, offset, limit uint64) ([]*model.Planet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Query", ctx, query, offset, limit)
	ret0, _ := ret[0].([]*model.Planet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockPlanetStoreMockRecorder) Query(ctx, query, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockPlanetStore)(nil).Query), ctx, query, offset, limit)
}

// Update mocks base method.
func (m_2 *MockPlanetStore) Update(ctx context.Context, m *model.Planet) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "Update", ctx, m)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockPlanetStoreMockRecorder) Update(ctx, m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPlanetStore)(nil).Update), ctx, m)
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