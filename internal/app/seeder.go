package app

import (
	"context"

	"github.com/kamva/tracer"
	"t3.org/t3/internal/input"
)

func (a *appCore) Seed(ctx context.Context, count int32) error {
	for i := int32(0); i < count; i++ {
		in := input.RandomCreatTicket()
		if _, err := a.createTicket(ctx, in, false); err != nil {
			return tracer.Trace(err)
		}
	}

	return nil
}

var _ Seeder = &appCore{}
