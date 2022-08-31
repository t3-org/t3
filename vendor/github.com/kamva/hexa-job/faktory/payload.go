package hexafaktory

import (
	"github.com/kamva/gutil"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
)

type mapPayload struct {
	val map[string]interface{}
}

func (p *mapPayload) Decode(payload interface{}) error {
	return tracer.Trace(gutil.MapToStruct(p.val, payload))
}

func newMapPayload(val map[string]interface{}) hjob.Payload {
	return &mapPayload{val: val}
}

var _ hjob.Payload = &mapPayload{}
