package hevent

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/kamva/tracer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Encoder encode and decode the event payload.
type Encoder interface {
	// Name returns the Encoder name.
	Name() string
	Encode(interface{}) ([]byte, error)
	Decoder([]byte) Decoder
}

// Decoder is event payload decoder.
type Decoder interface {
	// Decode decodes payload to the provided value.
	Decode(val interface{}) error
}

// protobufEncoder is protobuf implementation of the Encoder
type protobufEncoder struct{}
type protobufDecoder struct {
	b []byte
}

// jsonEncoder is json implementation of the Encoder.
type jsonEncoder struct{}

const (
	jsonEncoderName     = "json"
	protobufEncoderName = "protobuf"
)

var encoders = map[string]Encoder{
	jsonEncoderName:     NewJsonEncoder(),
	protobufEncoderName: NewProtobufEncoder(),
}

var (
	protobufTypeErr = errors.New("the provided value is not protobuf message")
)

func (m jsonEncoder) Name() string {
	return jsonEncoderName
}

func (m jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m jsonEncoder) Decoder(buf []byte) Decoder {
	return json.NewDecoder(bytes.NewReader(buf))
}

func (m protobufEncoder) Name() string {
	return protobufEncoderName
}

func (m protobufEncoder) Encode(v interface{}) ([]byte, error) {
	// I think we can use message v2 here.
	pb, ok := v.(proto.Message)
	if !ok {
		return nil, tracer.Trace(protobufTypeErr)
	}
	return protojson.Marshal(pb)
}

func (m protobufEncoder) Decoder(buf []byte) Decoder {
	return &protobufDecoder{b: buf}
}

func (p *protobufDecoder) Decode(val interface{}) error {
	return protojson.Unmarshal(p.b, val.(proto.Message))
}

// NewJsonEncoder returns new instance of the json encoder.
func NewJsonEncoder() Encoder {
	return &jsonEncoder{}
}

// NewProtobufEncoder returns new instance of the protobuf encoder.
func NewProtobufEncoder() Encoder {
	return &protobufEncoder{}
}

var _ Encoder = jsonEncoder{}
var _ Encoder = protobufEncoder{}
var _ Decoder = &protobufDecoder{}
