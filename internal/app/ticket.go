package app

import (
	"context"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/pagination"
	"github.com/kamva/tracer"
	"k8s.io/apimachinery/pkg/labels"
	"t3.org/t3/internal/dto"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/helpers"
	"t3.org/t3/internal/input"
	"t3.org/t3/internal/model"
	"t3.org/t3/internal/service/channel"
)

func (a *appCore) TicketEditionUrl(ctx context.Context, key, val string) (string, error) {
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
	globalFingerprints, err := in.GlobalFingerprints()
	if err != nil {
		return nil, tracer.Trace(err)
	}

	l, err := a.store.Ticket().GetAllByGlobalFingerprint(ctx, globalFingerprints)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	tickets := model.TicketsMapByGlobalFingerprint(l...)

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
		t := tickets[helpers.GlobalFingerprint(val.StartedAt, val.Fingerprint)]
		if t == nil {
			creationInput := input.CreateTicket(*val)
			t, err = a.createTicket(ctx, &creationInput, true)
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

	return a.createTicket(ctx, in, true)
}

func (a *appCore) createTicket(ctx context.Context, in *input.CreateTicket, dispatch bool) (*model.Ticket, error) {
	var ticket model.Ticket
	if err := ticket.Create(in); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := a.store.Ticket().Create(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}

	if dispatch {
		if err := a.dispatcher.Dispatch(ctx, &channel.DispatchInput{Ticket: &ticket}); err != nil {
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
	// TODO: In case the ticket is resolved and we get a firing input from the ticketGenerator,
	// we'll dispatch the ticket again, but because we don't change the spam status, it just send
	// a resolved ticket again, we can either do not send ticket in this situation or just create a
	// copy of it with firing=true status and dispatch it instead of the real ticket.
	isChangedFiring := in.IsFiring != nil && *in.IsFiring != t.IsFiring

	if err := t.Patch(in); err != nil {
		return nil, tracer.Trace(err)
	}
	if err := a.store.Ticket().Update(ctx, t); err != nil {
		return nil, tracer.Trace(err)
	}

	dispatchIn := &channel.DispatchInput{Ticket: t}
	if in.IsFromTicketGenerator {
		dispatchIn.IsFiringOnTicketGenerator = in.IsFiring
	}

	if isChangedFiring || in.IsFromTicketGenerator {
		if err := a.dispatcher.Dispatch(ctx, dispatchIn); err != nil {
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
	q, err := labels.Parse(query)
	if err != nil {
		return nil, apperr.ErrInvalidQuery.SetData(hexa.Map{"query": query})
	}
	count, err := a.store.Ticket().Count(ctx, q)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	p := pagination.New(page, perPage, count)
	p.Items, err = a.store.Ticket().Query(ctx, q, uint64(p.Offset()), uint64(p.Limit()))

	// TODO: Return ([]*model.Ticket,*pagination.Pages,error) as return result.
	return p, tracer.Trace(err)
}

var _ TicketService = &appCore{}
