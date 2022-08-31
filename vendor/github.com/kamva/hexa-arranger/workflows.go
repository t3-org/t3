package arranger

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type Workflows struct {
}

// Ctx converts workflow context to hexa context.
func (ac Workflows) Ctx(ctx workflow.Context) context.Context {
	return HexaCtxFromCadenceCtx(ctx)
}
