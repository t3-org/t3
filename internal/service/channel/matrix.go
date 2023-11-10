package channel

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
	"t3.org/t3/pkg/md"
)

type MatrixChannel struct {
	o       MatrixChannelOpts
	cli     *mautrix.Client
	kvStore KVStore
	md      *md.Markdown
}

type MatrixChannelOpts struct {
	// Set keyPrefix with "_" to set that label as an internal label.
	KeyPrefix     string
	OkEmoji       string
	CommandPrefix string
}

func NewMatrixChannel(cli *mautrix.Client, kv KVStore, md *md.Markdown, o MatrixChannelOpts) Channel {
	return &MatrixChannel{cli: cli, kvStore: kv, o: o, md: md}
}

func (m *MatrixChannel) Options() MatrixChannelOpts {
	return m.o
}

func (m *MatrixChannel) Key(roomID id.RoomID) string {
	return fmt.Sprintf("%s:%s", m.o.KeyPrefix, roomID)
}

func (m *MatrixChannel) Client() *mautrix.Client {
	return m.cli
}

func (m *MatrixChannel) Firing(ctx context.Context, channelID string, t *model.Ticket) error {
	roomID := id.RoomID(channelID)
	res, err := m.sendText(roomID, "", fmt.Sprintf("__New Firing Ticket__ \n\n %s ", t.Markdown()))
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(m.kvStore.Set(ctx, t.ID, m.Key(roomID), string(res.EventID)))
}

func (m *MatrixChannel) Resolved(ctx context.Context, channelID string, t *model.Ticket) error {
	roomID := id.RoomID(channelID)
	eventID, err := m.kvStore.Val(ctx, t.ID, m.Key(roomID))
	if err != nil && !errors.Is(err, apperr.ErrTicketKVNotFound) {
		return tracer.Trace(err)
	}

	_, err = m.sendText(roomID, id.EventID(eventID), fmt.Sprintf("__Resolved Ticket__ \n\n %s ", t.Markdown()))
	return tracer.Trace(err)
}

func (m *MatrixChannel) Shutdown(ctx context.Context) error {
	c, ok := m.cli.Crypto.(*cryptohelper.CryptoHelper)
	if ok {
		return tracer.Trace(c.Close())
	}
	return nil
}

func (m *MatrixChannel) sendText(roomID id.RoomID, threadEventId id.EventID, msg string) (*mautrix.RespSendEvent, error) {
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

var _ Channel = &MatrixChannel{}
var _ hexa.Shutdownable = &MatrixChannel{}
