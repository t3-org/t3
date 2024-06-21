package channel

import (
	"context"
)

type KVStore interface {
	// Set sets the value and update it if it already existed.
	Set(ctx context.Context, ticketID string, key string, val string) error
	Val(ctx context.Context, ticketID string, key string) (string, error)
}

// Home is channel home. it's like both client/server for a channel.
// We use the Home to send messages to the channel server(like a client).
// And also use it to serve the channel's commands from the channel's server.
type Home interface {
	Dispatch(ctx context.Context, cfg any, in *DispatchInput) error
}
