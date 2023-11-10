package app

import (
	"context"
	"fmt"

	"t3.org/t3/internal/input"
	"t3.org/t3/internal/model"
)

func (a *appCore) callTicketWebhook(ctx context.Context, in *input.Channel, t *model.Ticket) error {
	ch, ok := a.channels[in.Name]
	if !ok {
		return fmt.Errorf("can not find channel with name: %s", in.Name)
	}

	if t.IsFiring {
		return ch.Firing(ctx, in.ID, t)
	}

	return ch.Resolved(ctx, in.ID, t)
}
