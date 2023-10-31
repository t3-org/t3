package app

import (
	"context"

	"github.com/kamva/hexa/pagination"
	"github.com/kamva/tracer"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
	"space.org/space/internal/model"
)

func (a *appCore) GetTicket(ctx context.Context, id int64) (*dto.Ticket, error) {
	return a.store.Ticket().Get(ctx, id)
}

func (a *appCore) GetTicketByCode(ctx context.Context, code string) (*dto.Ticket, error) {
	return a.store.Ticket().GetByCode(ctx, code)
}

func (a *appCore) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	var ticket model.Ticket
	if err := ticket.Create(in); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := a.store.Ticket().Create(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}
	return &ticket, nil
}

func (a *appCore) PatchTicket(ctx context.Context, id int64, in *input.UpdateTicket) (*dto.Ticket, error) {
	ticket, err := a.store.Ticket().Get(ctx, id)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	if err := ticket.Patch(in); err != nil {
		return nil, tracer.Trace(err)
	}
	if err := a.store.Ticket().Update(ctx, ticket); err != nil {
		return nil, tracer.Trace(err)
	}
	return ticket, nil
}

func (a *appCore) DeleteTicket(ctx context.Context, id int64) error {
	ticket, err := a.store.Ticket().Get(ctx, id)
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(a.store.Ticket().Delete(ctx, ticket))
}

// TODO: Return ([]*model.Ticket,*pagination.Pages,error) as return result.

//nolint:revive
func (a *appCore) QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	// TODO: implement me
	panic("implement me")
}

var _ TicketService = &appCore{}
