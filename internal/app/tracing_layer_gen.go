// Code generated by "make app-layers"
// DO NOT EDIT
package app

import (
	"context"
	"reflect"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/pagination"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"space.org/space/internal/dto"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/input"
)

type tracingLayer struct {
	t    trace.Tracer
	next App
}

func (a *tracingLayer) HealthIdentifier() string {
	return a.next.HealthIdentifier()
}
func (a *tracingLayer) LivenessStatus(ctx context.Context) hexa.LivenessStatus {
	if ctx == nil {
		return a.next.LivenessStatus(ctx)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "LivenessStatus")
	defer span.End()

	r1 := a.next.LivenessStatus(ctx)
	span.SetStatus(codes.Ok, "")

	return r1
}
func (a *tracingLayer) ReadinessStatus(ctx context.Context) hexa.ReadinessStatus {
	if ctx == nil {
		return a.next.ReadinessStatus(ctx)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "ReadinessStatus")
	defer span.End()

	r1 := a.next.ReadinessStatus(ctx)
	span.SetStatus(codes.Ok, "")

	return r1
}
func (a *tracingLayer) HealthStatus(ctx context.Context) hexa.HealthStatus {
	if ctx == nil {
		return a.next.HealthStatus(ctx)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "HealthStatus")
	defer span.End()

	r1 := a.next.HealthStatus(ctx)
	span.SetStatus(codes.Ok, "")

	return r1
}
func (a *tracingLayer) GetPlanet(ctx context.Context, id int64) (*dto.Planet, error) {
	if ctx == nil {
		return a.next.GetPlanet(ctx, id)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "GetPlanet")
	defer span.End()

	r1, r2 := a.next.GetPlanet(ctx, id)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) GetPlanetByCode(ctx context.Context, code string) (*dto.Planet, error) {
	if ctx == nil {
		return a.next.GetPlanetByCode(ctx, code)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "GetPlanetByCode")
	defer span.End()

	r1, r2 := a.next.GetPlanetByCode(ctx, code)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) CreatePlanet(ctx context.Context, in input.CreatePlanet) (*dto.Planet, error) {
	if ctx == nil {
		return a.next.CreatePlanet(ctx, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "CreatePlanet")
	defer span.End()

	r1, r2 := a.next.CreatePlanet(ctx, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) UpdatePlanet(ctx context.Context, id int64, in input.UpdatePlanet) (*dto.Planet, error) {
	if ctx == nil {
		return a.next.UpdatePlanet(ctx, id, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "UpdatePlanet")
	defer span.End()

	r1, r2 := a.next.UpdatePlanet(ctx, id, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) DeletePlanet(ctx context.Context, id int64) error {
	if ctx == nil {
		return a.next.DeletePlanet(ctx, id)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "DeletePlanet")
	defer span.End()

	r1 := a.next.DeletePlanet(ctx, id)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (a *tracingLayer) QueryPlanets(ctx context.Context, query string, page int, perPage int) (*pagination.Pages, error) {
	if ctx == nil {
		return a.next.QueryPlanets(ctx, query, page, perPage)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "QueryPlanets")
	defer span.End()

	r1, r2 := a.next.QueryPlanets(ctx, query, page, perPage)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}

func NewTracingLayer(tp trace.TracerProvider, next App) App {
	return &tracingLayer{
		t:    tp.Tracer(reflect.TypeOf(tracingLayer{}).PkgPath()),
		next: next,
	}
}
