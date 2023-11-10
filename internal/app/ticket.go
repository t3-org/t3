package app

import (
	"context"
	"strings"

	"github.com/kamva/hexa/pagination"
	"github.com/kamva/tracer"
	"t3.org/t3/internal/dto"
	"t3.org/t3/internal/input"
	"t3.org/t3/internal/model"
)

func (a *appCore) EditTicketUrlByKey(ctx context.Context, key, val string) (string, error) {
	ticket, err := a.store.Ticket().FirstByTicketLabel(ctx, key, val)
	if err != nil {
		return "", tracer.Trace(err)
	}
	r := strings.NewReplacer("{id}", ticket.ID)
	return r.Replace(a.cfg.UI.EditTicketURL), nil
}

func (a *appCore) GetTicketByKey(ctx context.Context, key, val string) (*dto.Ticket, error) {
	return a.store.Ticket().FirstByTicketLabel(ctx, key, val)
}

func (a *appCore) GetTicket(ctx context.Context, id string) (*dto.Ticket, error) {
	return a.store.Ticket().Get(ctx, id)
}

func (a *appCore) UpsertTickets(ctx context.Context, in *input.BatchUpsertTickets) ([]*dto.Ticket, error) {
	l, err := a.store.Ticket().GetAllByFingerprint(ctx, in.Fingerprints())
	if err != nil {
		return nil, tracer.Trace(err)
	}
	tickets := model.TicketsMapByFingerprint(l...)

	// Validate data
	for _, in := range in.Tickets {
		// mark this input as the creation input instead of patch input if ticket is nil.
		in.SetIsCreation(tickets[in.Fingerprint] == nil)
	}
	if err := v(ctx, in); err != nil {
		return nil, tracer.Trace(err)
	}

	// upsert
	finalTickets := make([]*dto.Ticket, len(in.Tickets))
	for i, val := range in.Tickets {
		t := tickets[val.Fingerprint]
		if t == nil {
			creationInput := input.CreateTicket(*val)
			t, err = a.createTicket(ctx, &creationInput)
			if err != nil {
				return nil, tracer.Trace(err)
			}
		} else if _, err := a.patchTicket(ctx, t, val); err != nil {
			return nil, tracer.Trace(err)
		}
		finalTickets[i] = t
	}

	return finalTickets, nil
}

func (a *appCore) CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error) {
	if err := v(ctx, in); err != nil {
		return nil, tracer.Trace(err)
	}

	return a.createTicket(ctx, in)
}

func (a *appCore) createTicket(ctx context.Context, in *input.CreateTicket) (*model.Ticket, error) {
	var ticket model.Ticket
	if err := ticket.Create(in); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := a.store.Ticket().Create(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}

	if in.Channel != nil {
		if err := a.callTicketWebhook(ctx, in.Channel, &ticket); err != nil {
			return nil, err
		}
	}

	return &ticket, nil
}

func (a *appCore) PatchTicket(ctx context.Context, id string, in *input.PatchTicket) (*dto.Ticket, error) {
	ticket, err := a.store.Ticket().Get(ctx, id)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return a.patchTicket(ctx, ticket, in)
}

func (a *appCore) PatchTicketByLabel(ctx context.Context, key, val string, in *input.PatchTicket) (*dto.Ticket, error) {
	ticket, err := a.store.Ticket().FirstByTicketLabel(ctx, key, val)
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

	if in.Channel != nil && in.IsFiring != nil {
		if err := a.callTicketWebhook(ctx, in.Channel, t); err != nil {
			return nil, tracer.Trace(err)
		}
	}

	return t, nil
}

func (a *appCore) DeleteTicket(ctx context.Context, id string) error {
	ticket, err := a.store.Ticket().Get(ctx, id)
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(a.store.Ticket().Delete(ctx, ticket))
}

func (a *appCore) QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	count, err := a.store.Ticket().Count(ctx, query)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	p := pagination.New(page, perPage, count)
	p.Items, err = a.store.Ticket().Query(ctx, query, uint64(p.Offset()), uint64(p.Limit()))

	// TODO: Return ([]*model.Ticket,*pagination.Pages,error) as return result.
	return p, tracer.Trace(err)
}

var _ TicketService = &appCore{}
