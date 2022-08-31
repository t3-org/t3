package arranger

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

// HexaCtxFromCadenceCtx extracts hexa context from Cadence context.
func HexaCtxFromCadenceCtx(ctx workflow.Context) context.Context {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		return nil
	}
	return hexaCtx.(context.Context)
}
