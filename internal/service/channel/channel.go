package channel

import (
	"context"

	"space.org/space/internal/model"
)

type KVStore interface {
	Set(ctx context.Context, ticketID int64, key string, val string) error
	Val(ctx context.Context, ticketID int64, key string) (string, error)
}

type Channel interface {
	Firing(ctx context.Context, channelID string, ticket *model.Ticket) error
	Resolved(ctx context.Context, channelID string, ticket *model.Ticket) error
}
