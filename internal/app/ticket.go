package app

import (
	"context"
	"strconv"
	"strings"

	"github.com/kamva/hexa/pagination"
	"github.com/kamva/tracer"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
	"space.org/space/internal/model"
)

func (a *appCore) EditTicketUrlByKey(ctx context.Context, key, val string) (string, error) {
	ticket, err := a.store.Ticket().GetByTicketKeyVal(ctx, key, val)
	if err != nil {
		return "", tracer.Trace(err)
	}
	r := strings.NewReplacer("{id}", strconv.FormatInt(ticket.ID, 10))
	return r.Replace(a.cfg.UI.EditTicketURL), nil
}

func (a *appCore) GetTicketByKey(ctx context.Context, key, val string) (*dto.Ticket, error) {
	return a.store.Ticket().GetByTicketKeyVal(ctx, key, val)
}

func (a *appCore) GetTicket(ctx context.Context, id int64) (*dto.Ticket, error) {
	return a.store.Ticket().Get(ctx, id)
}

func (a *appCore) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	var ticket model.Ticket
	if err := ticket.Create(in); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := a.store.Ticket().Create(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}

	if in.Webhook != nil {
		if err := a.callTicketWebhook(ctx, in.Webhook, &ticket); err != nil {
			return nil, err
		}
	}

	return &ticket, nil
}

func (a *appCore) PatchTicket(ctx context.Context, id int64, in *input.PatchTicket) (*dto.Ticket, error) {
	ticket, err := a.store.Ticket().Get(ctx, id)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return a.patchTicket(ctx, ticket, in)
}

func (a *appCore) PatchTicketByKey(ctx context.Context, key, val string, in *input.PatchTicket) (*dto.Ticket, error) {
	ticket, err := a.store.Ticket().GetByTicketKeyVal(ctx, key, val)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return a.patchTicket(ctx, ticket, in)
}

func (a *appCore) patchTicket(ctx context.Context, t *model.Ticket, in *input.PatchTicket) (*dto.Ticket, error) {
	if err := t.Patch(in); err != nil {
		return nil, tracer.Trace(err)
	}
	if err := a.store.Ticket().Update(ctx, t); err != nil {
		return nil, tracer.Trace(err)
	}

	if in.Webhook != nil && in.IsFiring != nil {
		if err := a.callTicketWebhook(ctx, in.Webhook, t); err != nil {
			return nil, err
		}
	}

	return t, nil
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
