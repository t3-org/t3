package hsynq

import (
	"encoding/json"

	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
)

// Headers are job headers that we can set on each task.
type Headers map[string][]byte

// Transformer transforms job to bytes and convert
// bytes to hexa context and job payloads again.
type Transformer interface {
	BytesFromJob(headers Headers, payload interface{}) ([]byte, error)
	PayloadFromBytes(b []byte) (Headers, hjob.Payload, error)
}

type JsonJob struct {
	Headers map[string][]byte
	Payload interface{} `json:"payload"`
}

type UnmarshalJsonJob struct {
	Headers map[string][]byte
	Payload json.RawMessage `json:"payload"`
}

type jsonTransformer struct{}

func NewJsonTransformer() Transformer {
	return &jsonTransformer{}
}

func (t *jsonTransformer) BytesFromJob(headers Headers, payload interface{}) ([]byte, error) {
	job := JsonJob{
		Headers: headers,
		Payload: payload,
	}

	b, err := json.Marshal(&job)
	return b, tracer.Trace(err)
}

func (t *jsonTransformer) PayloadFromBytes(b []byte) (Headers, hjob.Payload, error) {
	var payload UnmarshalJsonJob
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, nil, tracer.Trace(err)
	}
	return payload.Headers, hjob.NewJsonPayload(payload.Payload), nil
}

var _ Transformer = &jsonTransformer{}
