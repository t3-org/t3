package matrix

import (
	"context"
	"fmt"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/config"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/input"
	"space.org/space/internal/registry/services"
	"space.org/space/internal/service/channel"
)

func registerTicketCommands(r *Router, res *ticketResource) {
	r.Register("new", res.NewTicketLink, "get a link to the ticket creation form in the UI.") // non-thread command
	r.Register("edit", res.EditTicketLink, "Get a link to the ticket edition form in the UI.")
	r.Register("ticket", res.GetTicket, "Get the ticket.")

	r.Register("seen", res.SetSeen, "mark ticket as seen. e.g., `!!seen {minutes(default: 0)`")
	r.Register("spam", res.SetSpam, "set the spam flag on the ticket. e.g., `!!spam`, `!!spam false`, `!!spam true(default)`")
	r.Register("resolved", res.SetAsResolved, "set a ticket as resolved. e.g., `!!resolved {minutes(default: 0)}`")
	r.Register("firing", res.SetAsFiring, "set a ticket as firing.")
	r.Register("level", res.SetLevel, "set level of a ticket. e.g., `!!level {level: low,medium or high}`")
	r.Register("description", res.SetDescription, "set description on a ticket. It'll remove the previous description content if it's not empty. e.g., `!!description {msg}`")
}

type ticketResource struct {
	*Resource
	cfg *config.Config
	mx  *channel.MatrixChannel
	app app.App
}

func newTicketResource(s services.Services, app app.App) *ticketResource {
	return &ticketResource{
		Resource: NewResource(s),
		cfg:      s.Config(),
		mx:       s.Matrix(),
		app:      app,
	}
}

func (r *ticketResource) patch(ctx context.Context, cmd *Command, in *input.PatchTicket) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	_, err := r.app.PatchTicketByLabel(ctx, r.mx.Key(cmd.Event.RoomID), threadID, in)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendOKReaction(cmd.Event)
}

func (r *ticketResource) NewTicketLink(_ context.Context, cmd *Command) error {
	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("[Create a new one here](%s)", r.cfg.UI.NewTicketUrl))
}

func (r *ticketResource) EditTicketLink(ctx context.Context, cmd *Command) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	url, err := r.app.EditTicketUrlByKey(ctx, r.mx.Key(cmd.Event.RoomID), threadID)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("[Edit it here](%s)", url))
}

func (r *ticketResource) GetTicket(ctx context.Context, cmd *Command) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	t, err := r.app.GetTicketByKey(ctx, r.mx.Key(cmd.Event.RoomID), threadID)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("__Ticket__ \n\n %s ", t.Markdown()))
}

func (r *ticketResource) SetSeen(ctx context.Context, cmd *Command) error {
	at := time.Now().Add(time.Minute * time.Duration(gutil.ParseInt(cmd.Params, 0))).UnixMilli()
	return r.patch(ctx, cmd, &input.PatchTicket{SeenAt: &at})
}

func (r *ticketResource) SetSpam(ctx context.Context, cmd *Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{IsSpam: gutil.NewBool(cmd.Params != "false")})
}

func (r *ticketResource) SetAsResolved(ctx context.Context, cmd *Command) error {
	at := time.Now().Add(time.Minute * time.Duration(gutil.ParseInt(cmd.Params, 0))).UnixMilli()
	return r.patch(ctx, cmd, &input.PatchTicket{
		IsFiring: gutil.NewBool(false),
		EndedAt:  &at,
	})
}

func (r *ticketResource) SetAsFiring(ctx context.Context, cmd *Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{IsFiring: gutil.NewBool(true)})
}

func (r *ticketResource) SetLevel(ctx context.Context, cmd *Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{Level: &cmd.Params})
}

func (r *ticketResource) SetDescription(ctx context.Context, cmd *Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{Description: &cmd.Params})
}
