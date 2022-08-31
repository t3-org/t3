package hevent

import (
	"context"
	"errors"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

const (
	HeaderKeyReplyChannel   = "_reply_channel"
	HeaderKeyPayloadEncoder = "_payload_encoder" // the message body.
)

type RawMessageConverter interface {
	EventToRaw(c context.Context, e *Event) (*RawMessage, error)
	// RawMsgToMessage converts the raw message to a message.
	// primary is the primary driver's message that its receiver will get.
	RawMsgToMessage(c context.Context, raw *RawMessage, primary interface{}) (context.Context, Message, error)
}

type rawMessageConverter struct {
	p hexa.ContextPropagator
	e Encoder
}

func NewRawMessageConverter(p hexa.ContextPropagator, e Encoder) RawMessageConverter {
	return &rawMessageConverter{
		p: p,
		e: e,
	}
}

func (m *rawMessageConverter) EventToRaw(ctx context.Context, event *Event) (*RawMessage, error) {
	payload, err := m.e.Encode(event.Payload)

	if err != nil {
		return nil, tracer.Trace(err)
	}

	headers, err := m.p.Inject(ctx)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	headers[HeaderKeyReplyChannel] = []byte(event.ReplyChannel)
	headers[HeaderKeyPayloadEncoder] = []byte(m.e.Name())

	return &RawMessage{
		Headers: headers,
		Payload: payload,
	}, err
}

func (m *rawMessageConverter) RawMsgToMessage(c context.Context, rawMsg *RawMessage, primary interface{}) (
	ctx context.Context, msg Message, err error) {

	ctx, err = m.p.Extract(c, rawMsg.Headers)
	if err != nil {
		err = tracer.Trace(err)
		return
	}

	encoderName := string(rawMsg.Headers[HeaderKeyPayloadEncoder])
	encoder, ok := encoders[encoderName]
	if !ok {
		err = errors.New("can not find message payload's encoder/decoder")
		return
	}

	msg = Message{
		Primary:       primary,
		Headers:       rawMsg.Headers,
		CorrelationId: hexa.CtxCorrelationId(ctx),
		ReplyChannel:  string(rawMsg.Headers[HeaderKeyReplyChannel]),
		Payload:       encoder.Decoder(rawMsg.Payload),
	}

	return
}

var _ RawMessageConverter = &rawMessageConverter{}
