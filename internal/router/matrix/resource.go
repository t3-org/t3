package matrix

import (
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	"t3.org/t3/internal/registry/services"
	"t3.org/t3/pkg/md"
)

type Resource struct {
	okEmoji string
	cli     *mautrix.Client
	md      *md.Markdown
}

func NewResource(s services.Services) *Resource {
	return &Resource{
		okEmoji: s.Matrix().Options().OkEmoji,
		cli:     s.Matrix().Client(),
		md:      s.Markdown(),
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
	return r.sendReaction(e, r.okEmoji)
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
