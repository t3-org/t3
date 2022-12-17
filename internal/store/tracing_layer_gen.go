// Code generated by "make app-layers"
// DO NOT EDIT
package store

import (
	"context"
	"reflect"

	"github.com/kamva/hexa"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/pkg/sqld"
)

type tracingLayerStore struct {
	next        model.Store
	PlanetStore model.PlanetStore
	SystemStore model.SystemStore
}

func (s *tracingLayerStore) DBLayer() model.Store {
	return s.next.DBLayer()
}
func (s *tracingLayerStore) Tx() *sqld.TxWrapper {
	return s.next.Tx()
}
func (s *tracingLayerStore) TruncateAllTables(ctx context.Context) error {
	return s.next.TruncateAllTables(ctx)
}
func (s *tracingLayerStore) System() model.SystemStore {
	return s.SystemStore
}
func (s *tracingLayerStore) Planet() model.PlanetStore {
	return s.PlanetStore
}
func (s *tracingLayerStore) HealthIdentifier() string {
	return s.next.HealthIdentifier()
}
func (s *tracingLayerStore) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	return s.next.LivenessStatus(ctx)
}
func (s *tracingLayerStore) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	return s.next.ReadinessStatus(ctx)
}
func (s *tracingLayerStore) HealthStatus(ctx context.Context) hexa.HealthStatus {
	return s.next.HealthStatus(ctx)
}

func NewTracingLayerStore(instrumentationPostfix string, tp trace.TracerProvider, next model.Store) model.Store {
	pkgPath := reflect.TypeOf(tracingLayerStore{}).PkgPath() + "." + instrumentationPostfix

	return &tracingLayerStore{
		next:        next,
		PlanetStore: &tracingLayerPlanetStore{t: tp.Tracer(pkgPath), next: next.Planet()},
		SystemStore: &tracingLayerSystemStore{t: tp.Tracer(pkgPath), next: next.System()},
	}
}

//--------------------------------
// Define subStores
//--------------------------------

type tracingLayerPlanetStore struct {
	t    trace.Tracer
	next model.PlanetStore
}

func (s *tracingLayerPlanetStore) Get(ctx context.Context, id int64) (*model.Planet, error) {
	if ctx == nil {
		return s.next.Get(ctx, id)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Get")
	defer span.End()

	r1, r2 := s.next.Get(ctx, id)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (s *tracingLayerPlanetStore) GetByCode(ctx context.Context, code string) (*model.Planet, error) {
	if ctx == nil {
		return s.next.GetByCode(ctx, code)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.GetByCode")
	defer span.End()

	r1, r2 := s.next.GetByCode(ctx, code)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (s *tracingLayerPlanetStore) Create(ctx context.Context, m *model.Planet) error {
	if ctx == nil {
		return s.next.Create(ctx, m)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Create")
	defer span.End()

	r1 := s.next.Create(ctx, m)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (s *tracingLayerPlanetStore) Update(ctx context.Context, m *model.Planet) error {
	if ctx == nil {
		return s.next.Update(ctx, m)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Update")
	defer span.End()

	r1 := s.next.Update(ctx, m)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (s *tracingLayerPlanetStore) Delete(ctx context.Context, m *model.Planet) error {
	if ctx == nil {
		return s.next.Delete(ctx, m)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Delete")
	defer span.End()

	r1 := s.next.Delete(ctx, m)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (s *tracingLayerPlanetStore) Count(ctx context.Context, query string) (int, error) {
	if ctx == nil {
		return s.next.Count(ctx, query)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Count")
	defer span.End()

	r1, r2 := s.next.Count(ctx, query)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (s *tracingLayerPlanetStore) Query(ctx context.Context, query string, offset uint64, limit uint64) ([]*model.Planet, error) {
	if ctx == nil {
		return s.next.Query(ctx, query, offset, limit)
	}

	ctx, span := s.t.Start(ctx, "PlanetStore.Query")
	defer span.End()

	r1, r2 := s.next.Query(ctx, query, offset, limit)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}

type tracingLayerSystemStore struct {
	t    trace.Tracer
	next model.SystemStore
}

func (s *tracingLayerSystemStore) GetByName(ctx context.Context, name string) (*model.System, error) {
	if ctx == nil {
		return s.next.GetByName(ctx, name)
	}

	ctx, span := s.t.Start(ctx, "SystemStore.GetByName")
	defer span.End()

	r1, r2 := s.next.GetByName(ctx, name)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (s *tracingLayerSystemStore) Save(ctx context.Context, system *model.System) error {
	if ctx == nil {
		return s.next.Save(ctx, system)
	}

	ctx, span := s.t.Start(ctx, "SystemStore.Save")
	defer span.End()

	r1 := s.next.Save(ctx, system)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (s *tracingLayerSystemStore) Delete(ctx context.Context, name string) error {
	if ctx == nil {
		return s.next.Delete(ctx, name)
	}

	ctx, span := s.t.Start(ctx, "SystemStore.Delete")
	defer span.End()

	r1 := s.next.Delete(ctx, name)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
