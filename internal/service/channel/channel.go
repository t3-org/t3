package channel

import (
	"context"

	"space.org/space/internal/model"
)

type KVStore interface {
	// Set sets the value and update it if it already existed.
	Set(ctx context.Context, ticketID string, key string, val string) error
	Val(ctx context.Context, ticketID string, key string) (string, error)
}

type Channel interface {
	Firing(ctx context.Context, channelID string, ticket *model.Ticket) error
	Resolved(ctx context.Context, channelID string, ticket *model.Ticket) error
}
