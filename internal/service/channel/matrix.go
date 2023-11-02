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
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
)

type MatrixChannel struct {
	cli       *mautrix.Client
	kvStore   KVStore
	keyPrefix string // Prefix of the key that we use as the store's keys.
}

func NewMatrixChannel(cli *mautrix.Client, kv KVStore, keyPrefix string) Channel {
	return &MatrixChannel{cli: cli, kvStore: kv, keyPrefix: keyPrefix}
}

func (m *MatrixChannel) key(roomID id.RoomID) string {
	return fmt.Sprintf("%s:%s", m.keyPrefix, roomID)
}

func (m *MatrixChannel) Client() *mautrix.Client {
	return m.cli
}

func (m *MatrixChannel) Firing(ctx context.Context, channelID string, t *model.Ticket) error {
	roomID := id.RoomID(channelID)
	res, err := m.cli.SendText(roomID, fmt.Sprintf("New Firing Ticket: %+v", t))
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(m.kvStore.Set(ctx, t.ID, m.key(roomID), string(res.EventID)))
}

func (m *MatrixChannel) Resolved(ctx context.Context, channelID string, t *model.Ticket) error {
	roomID := id.RoomID(channelID)
	eventID, err := m.kvStore.Val(ctx, t.ID, m.key(roomID))
	if err != nil && !errors.Is(err, apperr.ErrTicketKVNotFound) {
		return tracer.Trace(err)
	}

	content := &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    fmt.Sprintf("New resolved Ticket: %+v", t),
	}

	if eventID != "" { // if there's a eventID for the thread of this ticket on this channel.
		content.RelatesTo = &event.RelatesTo{Type: event.RelThread, EventID: id.EventID(eventID)}
	}

	_, err = m.cli.SendMessageEvent(roomID, event.EventMessage, content)
	return tracer.Trace(err)
}

func (m *MatrixChannel) Shutdown(ctx context.Context) error {
	c, ok := m.cli.Crypto.(*cryptohelper.CryptoHelper)
	if ok {
		return tracer.Trace(c.Close())
	}
	return nil
}

var _ Channel = &MatrixChannel{}
var _ hexa.Shutdownable = &MatrixChannel{}
