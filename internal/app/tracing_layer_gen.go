// Code generated by "make layers"
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
func (a *tracingLayer) UpsertTickets(ctx context.Context, in *input.BatchUpsertTickets) ([]*dto.Ticket, error) {
	if ctx == nil {
		return a.next.UpsertTickets(ctx, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "UpsertTickets")
	defer span.End()

	r1, r2 := a.next.UpsertTickets(ctx, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) EditTicketUrlByKey(ctx context.Context, key string, val string) (string, error) {
	if ctx == nil {
		return a.next.EditTicketUrlByKey(ctx, key, val)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "EditTicketUrlByKey")
	defer span.End()

	r1, r2 := a.next.EditTicketUrlByKey(ctx, key, val)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) GetTicketByKey(ctx context.Context, key string, val string) (*dto.Ticket, error) {
	if ctx == nil {
		return a.next.GetTicketByKey(ctx, key, val)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "GetTicketByKey")
	defer span.End()

	r1, r2 := a.next.GetTicketByKey(ctx, key, val)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) GetTicket(ctx context.Context, id string) (*dto.Ticket, error) {
	if ctx == nil {
		return a.next.GetTicket(ctx, id)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "GetTicket")
	defer span.End()

	r1, r2 := a.next.GetTicket(ctx, id)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	if ctx == nil {
		return a.next.CreateTicket(ctx, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "CreateTicket")
	defer span.End()

	r1, r2 := a.next.CreateTicket(ctx, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) PatchTicket(ctx context.Context, id string, in *input.PatchTicket) (*dto.Ticket, error) {
	if ctx == nil {
		return a.next.PatchTicket(ctx, id, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "PatchTicket")
	defer span.End()

	r1, r2 := a.next.PatchTicket(ctx, id, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) PatchTicketByLabel(ctx context.Context, key string, val string, in *input.PatchTicket) (*dto.Ticket, error) {
	if ctx == nil {
		return a.next.PatchTicketByLabel(ctx, key, val, in)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "PatchTicketByLabel")
	defer span.End()

	r1, r2 := a.next.PatchTicketByLabel(ctx, key, val, in)

	if apperr.IsInternalErr(r2) {
		span.RecordError(r2)
		span.SetStatus(codes.Error, r2.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1, r2
}
func (a *tracingLayer) DeleteTicket(ctx context.Context, id string) error {
	if ctx == nil {
		return a.next.DeleteTicket(ctx, id)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "DeleteTicket")
	defer span.End()

	r1 := a.next.DeleteTicket(ctx, id)

	if apperr.IsInternalErr(r1) {
		span.RecordError(r1)
		span.SetStatus(codes.Error, r1.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return r1
}
func (a *tracingLayer) QueryTickets(ctx context.Context, query string, page int, perPage int) (*pagination.Pages, error) {
	if ctx == nil {
		return a.next.QueryTickets(ctx, query, page, perPage)
	}

	var span trace.Span
	ctx, span = a.t.Start(ctx, "QueryTickets")
	defer span.End()

	r1, r2 := a.next.QueryTickets(ctx, query, page, perPage)

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
