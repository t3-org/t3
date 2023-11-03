package matrix

import (
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Resource struct {
	okEmoji string
	cli     *mautrix.Client
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
