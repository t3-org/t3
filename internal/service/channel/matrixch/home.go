package matrixch

import (
	"context"
	"errors"
	"fmt"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/model"
	"t3.org/t3/internal/service/channel"
	"t3.org/t3/pkg/md"
)

type ChannelOptions struct {
	RoomID string `yaml:"room_id"`
}

type Home struct {
	*Server
	opts    *HomeOpts
	cli     *mautrix.Client
	kvStore channel.KVStore
	md      *md.Markdown
}

func New(srv *Server, o *HomeOpts) channel.Home {
	return &Home{
		opts:    o,
		Server:  srv,
		cli:     o.Client,
		kvStore: o.KVStore,
		md:      o.MarkDown,
	}
}
func (m *Home) Dispatch(ctx context.Context, conf any, t *model.Ticket) error {
	opts := conf.(ChannelOptions)
	title := fmt.Sprintf("__Firing Ticket__ \n\n %s ", t.Markdown())
	if !t.IsFiring {
		title = fmt.Sprintf("__Resolved Ticket__ \n\n %s ", t.Markdown())
	}

	return tracer.Trace(m.dispatch(ctx, opts.RoomID, t, title))
}

func (m *Home) dispatch(ctx context.Context, rID string, t *model.Ticket, title string) error {
	roomID := id.RoomID(rID)
	eventID, err := m.kvStore.Val(ctx, t.ID, m.opts.Key(roomID))
	if err != nil && !errors.Is(err, apperr.ErrTicketKVNotFound) {
		return tracer.Trace(err)
	}

	res, err := m.sendText(roomID, id.EventID(eventID), title)
	if err != nil {
		return tracer.Trace(err)
	}

	if eventID == "" {
		return tracer.Trace(m.kvStore.Set(ctx, t.ID, m.opts.Key(roomID), string(res.EventID)))
	}

	return nil
}

func (m *Home) Shutdown(ctx context.Context) error {
	if err := m.Server.doShutdown(ctx); err != nil {
		return tracer.Trace(err)
	}
	c, ok := m.cli.Crypto.(*cryptohelper.CryptoHelper)
	if ok {
		return tracer.Trace(c.Close())
	}
	return nil
}

func (m *Home) sendText(roomID id.RoomID, threadEventId id.EventID, msg string) (*mautrix.RespSendEvent, error) {
	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    msg,
	}

	rendered := m.md.RenderString(msg)
	if rendered != msg {
		content.Format = event.FormatHTML
		content.FormattedBody = rendered
	}

	if threadEventId != "" { // if there's an eventID for the thread of this ticket on this room.
		content.RelatesTo = &event.RelatesTo{Type: event.RelThread, EventID: threadEventId}
	}

	return m.cli.SendMessageEvent(roomID, event.EventMessage, content)
}

var _ channel.Home = &Home{}
var _ hexa.Bootable = &Home{}
var _ hexa.Runnable = &Home{}
var _ hexa.Shutdownable = &Home{}
