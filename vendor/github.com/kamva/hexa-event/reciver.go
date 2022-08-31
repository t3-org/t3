package hevent

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
)

// SubscriptionOptions contains options to subscribe to one or multiple channels.
type SubscriptionOptions struct {
	// Channel specify the channel name you will subscribe on.
	// Either Channel,Channels or ChannelsPattern are required when subscribing.
	Channel string

	// Channels contains name of channels which we want to subscribe.
	// Either Channel,Channels or ChannelsPattern are required when subscribing.
	Channels []string

	// ChannelsPattern is the pattern you will use to subscribe on all channels
	// which match with this pattern.
	// Either Channel,Channels or ChannelsPattern are required when subscribing.
	ChannelsPattern string

	// Handler is the event handler.
	Handler EventHandler

	// extra contains extra details for specific drivers(e.g for pulsar you can set extra consumer options here).
	extra []interface{}
}

// EventHandler handle events.
// pulsar and hestan implementations just log returned error, in kafka
// if you return error, it will push event to the retry or DLQ topic.
type EventHandler func(HandlerContext, Message, error) error

type Receiver interface {
	// Subscribe subscribe to the provided channel
	Subscribe(channel string, h EventHandler) error

	// SubscribeWithOptions subscribe by options.
	SubscribeWithOptions(*SubscriptionOptions) error

	hexa.Runnable     // to start receiving events.
	hexa.Shutdownable // to close connections and shutdown the server.
}

// HandlerContext is the context that pass to the message handler.
type HandlerContext interface {
	context.Context
	// Ack get the message and send ack.
	Ack()
	// Nack gets the message and send negative ack.
	Nack()
}

// RawMessage is the message sent by emitter,
// we will convert RawMessage to message and then
// pass it to the event handler.
// Note: Some event drivers (kafkabox & hafka) do
// not push the marshaled RawMessage as the event
// value, they send RawMessage's headers in the headers
// section and RawMessage's payload in the
// Payload section of the event, so if you want to define
// extra fields in addition to Headers and Payload in
// the RawMessage, please be careful.
type RawMessage struct {
	Headers map[string][]byte `json:"header,omitempty"`
	Payload []byte            `json:"payload"`
}

// Message is the message that provide to event handler.
type Message struct {
	// Primary is not the RawMessage. its the driver's
	// raw message.
	Primary interface{}

	Headers map[string][]byte

	CorrelationId string
	ReplyChannel  string
	Payload       Decoder
}

func (so *SubscriptionOptions) Validate() error {
	err := validation.ValidateStruct(so,
		validation.Field(&so.Channel, validation.Required.When(len(so.Channels) == 0 && so.ChannelsPattern == "")),
		validation.Field(&so.ChannelsPattern, validation.Required.When(so.Channel == "" && len(so.Channels) == 0)),
		validation.Field(
			&so.Channels,
			validation.Required.When(so.Channel == "" && so.ChannelsPattern == ""),
			validation.Each(validation.Required),
		),
	)
	return tracer.Trace(err)
}

// WithExtra add Extra data to the subscription options.
func (so *SubscriptionOptions) WithExtra(extra ...interface{}) *SubscriptionOptions {
	so.extra = append(so.extra, extra...)
	return so
}

// Extra returns the extra data of the subscription options.
func (so *SubscriptionOptions) Extra() []interface{} {
	return so.extra
}

func (m Message) Validate() error { // TODO: I think we should remove this method.
	if m.CorrelationId == "" {
		return tracer.Trace(errors.New("correlation-id is required the event"))
	}

	if m.Headers == nil || m.Payload == nil {
		return tracer.Trace(errors.New("message header and payload are required"))
	}
	return nil
}

func (e RawMessage) Validate() error {
	if e.Headers == nil {
		return tracer.Trace(errors.New("header is required in the raw event message"))
	}
	return nil
}

// NewSubscriptionOptions returns new instance of the subscription options.
func NewSubscriptionOptions(channel string, handler EventHandler) *SubscriptionOptions {
	return &SubscriptionOptions{
		Channel: channel,
		Handler: handler,
	}
}

// Assertion
var _ validation.Validatable = &SubscriptionOptions{}
var _ validation.Validatable = &Message{}
