package hcache

import (
	"bytes"
	"encoding/gob"

	"github.com/kamva/tracer"
	"github.com/tinylib/msgp/msgp"
	"github.com/vmihailenco/msgpack/v5"
)

type Marshaler func(val interface{}) ([]byte, error)

type Unmarshaler func(msg []byte, val interface{}) error

func MsgpackMarshal(val interface{}) ([]byte, error) {
	// The fast path (using generated code)
	if msgpVal, ok := val.(msgp.Marshaler); ok {
		return msgpVal.MarshalMsg(nil)
	}

	// The slow path
	return msgpack.Marshal(val)
}

func MsgpackUnmarshal(msg []byte, val interface{}) error {
	// The fast path (using generated code)
	if msgpVal, ok := val.(msgp.Unmarshaler); ok {
		_, err := msgpVal.UnmarshalMsg(msg)
		return tracer.Trace(err)
	}

	// The slow path
	return msgpack.Unmarshal(msg, &val)
}

func GobMarshal(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(v)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func GobUnmarshal(data []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(v)
}
