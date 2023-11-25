package command

import (
	"context"
	"fmt"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"gopkg.in/yaml.v2"
	"t3.org/t3/internal/app"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/input"
	"t3.org/t3/internal/service/channel/matrixch"
)

func registerTicketCommands(r *matrixch.Router, res *ticketResource) {
	r.Register("new", res.New, "Creaet a new ticket by providing its yaml data")                   // non-thread command
	r.Register("new_link", res.NewTicketLink, "get a link to the ticket creation form in the UI.") // non-thread command
	r.Register("patch", res.Patch, "Patch the ticket using yaml value passed as param of the command")
	r.Register("edit_link", res.EditTicketLink, "Get a link to the ticket edition form in the UI.")
	r.Register("ticket", res.GetTicket, "Get the ticket.")
	r.Register("ticket_yaml", res.GetTicketInYaml, "Get the ticket in yaml format.")

	r.Register("seen", res.SetSeen, "mark ticket as seen. e.g., `!!seen {minutes(default: 0)`")
	r.Register("spam", res.SetSpam, "set the spam flag on the ticket. e.g., `!!spam`, `!!spam false`, `!!spam true(default)`")
	r.Register("resolved", res.SetAsResolved, "set a ticket as resolved. e.g., `!!resolved {minutes(default: 0)}`")
	r.Register("firing", res.SetAsFiring, "set a ticket as firing.")
	r.Register("level", res.SetLevel, "set level of a ticket. e.g., `!!level {level: low,medium or high}`")
	r.Register("description", res.SetDescription, "set description on a ticket. It'll remove the previous description content if it's not empty. e.g., `!!description {msg}`")
}

type ticketResource struct {
	*Resource
	app app.App
}

func newTicketResource(res *Resource, app app.App) *ticketResource {
	return &ticketResource{
		Resource: res,
		app:      app,
	}
}

func (r *ticketResource) patch(ctx context.Context, cmd *matrixch.Command, in *input.PatchTicket) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	_, err := r.app.PatchTicketByLabel(ctx, r.o.Key(cmd.Event.RoomID), threadID, in)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendOKReaction(cmd.Event)
}

func (r *ticketResource) NewTicketLink(_ context.Context, cmd *matrixch.Command) error {
	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("[Create a new one here](%s)", r.o.UI.NewTicketUrl))
}

func (r *ticketResource) EditTicketLink(ctx context.Context, cmd *matrixch.Command) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	url, err := r.app.EditTicketUrlByKey(ctx, r.o.Key(cmd.Event.RoomID), threadID)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("[Edit it here](%s)", url))
}

func (r *ticketResource) GetTicket(ctx context.Context, cmd *matrixch.Command) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	t, err := r.app.GetTicketByKey(ctx, r.o.Key(cmd.Event.RoomID), threadID)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf("__Ticket__ \n\n %s ", t.Markdown()))
}

func (r *ticketResource) GetTicketInYaml(ctx context.Context, cmd *matrixch.Command) error {
	threadID := string(r.threadEventID(cmd.Event))
	if threadID == "" {
		return apperr.ErrTicketNotFound
	}

	t, err := r.app.GetTicketByKey(ctx, r.o.Key(cmd.Event.RoomID), threadID)
	if err != nil {
		return tracer.Trace(err)
	}
	res, err := yaml.Marshal(t)
	if err != nil {
		return tracer.Trace(err)
	}
	return r.SendTextWithSameRelation(cmd.Event, fmt.Sprintf(`<code><pre>%s</pre></code>`, string(res)))
}

func (r *ticketResource) New(ctx context.Context, cmd *matrixch.Command) error {
	var in input.CreateTicket
	if err := yaml.Unmarshal([]byte(cmd.Params), &in); err != nil {
		return tracer.Trace(err)
	}
	_, err := r.app.CreateTicket(ctx, &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return r.SendOKReaction(cmd.Event)
}

func (r *ticketResource) Patch(ctx context.Context, cmd *matrixch.Command) error {
	var in input.PatchTicket
	if err := yaml.Unmarshal([]byte(cmd.Params), &in); err != nil {
		return tracer.Trace(err)
	}
	return r.patch(ctx, cmd, &in)
}

func (r *ticketResource) SetSeen(ctx context.Context, cmd *matrixch.Command) error {
	at := time.Now().Add(time.Minute * time.Duration(gutil.ParseInt(cmd.Params, 0))).UnixMilli()
	return r.patch(ctx, cmd, &input.PatchTicket{SeenAt: &at})
}

func (r *ticketResource) SetSpam(ctx context.Context, cmd *matrixch.Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{IsSpam: gutil.NewBool(cmd.Params != "false")})
}

func (r *ticketResource) SetAsResolved(ctx context.Context, cmd *matrixch.Command) error {
	at := time.Now().Add(time.Minute * time.Duration(gutil.ParseInt(cmd.Params, 0))).UnixMilli()
	return r.patch(ctx, cmd, &input.PatchTicket{
		IsFiring: gutil.NewBool(false),
		EndedAt:  &at,
	})
}

func (r *ticketResource) SetAsFiring(ctx context.Context, cmd *matrixch.Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{IsFiring: gutil.NewBool(true)})
}

func (r *ticketResource) SetLevel(ctx context.Context, cmd *matrixch.Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{Severity: &cmd.Params})
}

func (r *ticketResource) SetDescription(ctx context.Context, cmd *matrixch.Command) error {
	return r.patch(ctx, cmd, &input.PatchTicket{Description: &cmd.Params})
}
