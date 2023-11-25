package command

import (
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"t3.org/t3/internal/service/channel/matrixch"
	"t3.org/t3/pkg/md"
)

type Resource struct {
	o   *matrixch.HomeOpts
	cli *mautrix.Client
	md  *md.Markdown
}

func NewResource(opts *matrixch.HomeOpts) *Resource {
	return &Resource{
		o:   opts,
		cli: opts.Client,
		md:  opts.MarkDown,
	}
}

func (r *Resource) threadEventID(e *event.Event) id.EventID {
	return e.Content.AsMessage().GetRelatesTo().GetThreadParent()
}

func (r *Resource) SendTextWithSameRelation(e *event.Event, text string) error {
	content := &event.MessageEventContent{
		MsgType:   event.MsgText,
		Body:      text,
		RelatesTo: e.Content.AsMessage().GetRelatesTo(),
	}

	rendered := r.md.RenderString(text)
	if rendered != text {
		content.Format = event.FormatHTML
		content.FormattedBody = rendered
	}
	_, err := r.cli.SendMessageEvent(e.RoomID, event.EventMessage, content)
	return tracer.Trace(err)
}

func (r *Resource) SendOKReaction(e *event.Event) error {
	return r.sendReaction(e, r.o.OkEmoji)
}

func (r *Resource) sendReaction(e *event.Event, key string) error {
	content := event.MessageEventContent{
		RelatesTo: &event.RelatesTo{
			Type:    event.RelAnnotation,
			EventID: e.ID,
			Key:     key, // key field is the emoji that we want to send.
		},
	}

	_, err := r.cli.SendMessageEvent(e.RoomID, event.EventReaction, content)
	return tracer.Trace(err)
}
