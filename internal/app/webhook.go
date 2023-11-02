package app

import (
	"context"
	"fmt"

	"space.org/space/internal/input"
	"space.org/space/internal/model"
)

func (a *appCore) callTicketWebhook(ctx context.Context, in *input.Webhook, t *model.Ticket) error {
	ch, ok := a.channels[in.Channel]
	if !ok {
		return fmt.Errorf("can not find channel with name: %s", in.Channel)
	}

	if t.IsFiring {
		return ch.Firing(ctx, in.ChannelID, t)
	}

	return ch.Resolved(ctx, in.ChannelID, t)
}
