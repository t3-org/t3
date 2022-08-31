package arranger

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/workflow"
)

// hexaCtxKey is the which we use set exported hexa context
// in cadence context propagator.
const hexaCtxKey = "_hexa_ctx"
const hexaKeys = "_hexa_ctx_keys"

// hexaContextPropagator propagate hexa context.
type hexaContextPropagator struct {
	p hexa.ContextPropagator
	// if set strict flag to true, you must set hexa
	// context on all calls to workflow,activity,....
	strict bool
}

func (h *hexaContextPropagator) Inject(ctx context.Context, hw workflow.HeaderWriter) error {
	m, err := h.p.Inject(ctx)
	if err != nil {
		return err
	}
	keys := make([]string, 0)

	for k, v := range m {
		hw.Set(k, payload(v))
		keys = append(keys, k)
	}
	hw.Set(hexaKeys, payload([]byte(strings.Join(keys, ","))))

	return nil
}

func (h *hexaContextPropagator) Extract(ctx context.Context, hr workflow.HeaderReader) (context.Context, error) {
	keysStr, ok := hr.Get(hexaKeys)
	if !ok {
		return nil, errors.New("can not find propagated hexa context key")
	}
	keys := strings.Split(string(keysStr.Data), ",")
	m := make(map[string][]byte)
	for _, k := range keys {
		val, ok := hr.Get(k)
		if !ok {
			err := fmt.Errorf("can not find %s hexa key context while strict mode is enabled", k)
			return nil, tracer.Trace(err)
		}
		m[k] = val.Data
	}

	return h.p.Extract(ctx, m)
}

func (h *hexaContextPropagator) InjectFromWorkflow(ctx workflow.Context, hw workflow.HeaderWriter) error {
	hexaCtx := ctx.Value(hexaCtxKey)
	if hexaCtx == nil {
		if h.strict {
			return errors.New("you must provide hexa context when strict mode is enabled")
		}
		return nil
	}
	m, err := h.p.Inject(hexaCtx.(context.Context))
	if err != nil {
		return err
	}

	// we should also keep the map's keys
	keys := make([]string, 0)
	for k, v := range m {
		hw.Set(k, payload(v))
		keys = append(keys, k)
	}
	hw.Set(hexaKeys, payload([]byte(strings.Join(keys, ","))))

	return nil
}

func (h *hexaContextPropagator) ExtractToWorkflow(ctx workflow.Context, hr workflow.HeaderReader) (workflow.Context, error) {
	keysStr, ok := hr.Get(hexaKeys)
	if !ok {
		return nil, errors.New("can not find propagated hexa context key")
	}
	keys := strings.Split(string(keysStr.Data), ",")

	m := make(map[string][]byte)
	for _, k := range keys {
		val, ok := hr.Get(k)
		if !ok {
			err := fmt.Errorf("can not find %s hexa key context while strict mode is enabled", k)
			return nil, tracer.Trace(err)
		}
		m[k] = val.Data
	}

	hexaCtx, err := h.p.Extract(context.Background(), m)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return workflow.WithValue(ctx, hexaCtxKey, hexaCtx), nil
}

// NewHexaContextPropagator returns new instance of hexa context propagator.
func NewHexaContextPropagator(p hexa.ContextPropagator) workflow.ContextPropagator {
	return &hexaContextPropagator{p: p}
}

func payload(data []byte) *common.Payload {
	return &common.Payload{Data: data}
}

var _ workflow.ContextPropagator = &hexaContextPropagator{}
