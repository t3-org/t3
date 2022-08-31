package gutil

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

// marshalNode defines function which marshal provided input.
// use it to create chain of nodes.
type marshalNode func(i interface{}) ([]byte, error)

// unmarshalNode defines function which unmarshal provided binary input.
// use it to create chain of nodes.
type unmarshalNode func(b []byte, o interface{}) error

// marshalProtobufNode is a node which marshal protobuf message
// if input type is not a protobuf message, it pass input
// to the next node.
func marshalProtobufNode(next marshalNode) marshalNode {
	return func(i interface{}) ([]byte, error) {
		m, ok := i.(proto.Message)
		if !ok {
			return next(i)
		}

		options := protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
			UseEnumNumbers:  true,
		}
		return options.Marshal(proto.MessageV2(m))
	}
}

// marshalJSONNode is a node which marshal input to binary data.
func marshalJSONNode(_ marshalNode) marshalNode {
	return func(i interface{}) ([]byte, error) {
		return json.Marshal(i)
	}
}

// marshalProtobufNode is a node which unmarshal provided
// binary data to the provided protobuf message.
// if type of output is not a protobuf message, it pass input
// to the next node.
func unmarshalProtobufNode(next unmarshalNode) unmarshalNode {
	return func(b []byte, o interface{}) error {
		m, ok := o.(proto.Message)
		if !ok {
			return next(b, o)
		}
		options := protojson.UnmarshalOptions{DiscardUnknown: true}
		return options.Unmarshal(b, proto.MessageV2(m))
	}
}

// unmarshalJSONNode unmarshal binary data to the provided
// output.
func unmarshalJSONNode(_ unmarshalNode) unmarshalNode {
	return func(b []byte, o interface{}) error {
		return json.Unmarshal(b, o)
	}
}

var (
	// marshalChain is the marshal chain
	marshalChain = marshalProtobufNode(marshalJSONNode(nil))
	// unmarshalChain is the unmarshal chain
	unmarshalChain = unmarshalProtobufNode(unmarshalJSONNode(nil))
)

func Marshal(in interface{}) ([]byte, error) {
	return marshalChain(in)
}

func Unmarshal(in []byte, out interface{}) error {
	return unmarshalChain(in, out)
}

// StructToMap convert the struct to a map by marshall and un-marshall
func StructToMap(input interface{}) map[string]interface{} {
	var m = map[string]interface{}{}

	encodedJson, _ := Marshal(input)
	_ = Unmarshal(encodedJson, &m)

	return m
}

// MapToStruct convert map to struct by json marshal and Unmarshal
func MapToStruct(m map[string]interface{}, s interface{}) error {
	encodedJson, err := Marshal(m)
	if err != nil {
		return err
	}

	return Unmarshal(encodedJson, s)
}

// UnmarshalStruct marshal provided input to json and then unmarshal json bytes to provided output.
// input and output can be regular struct, map or struct of protobuf message type.
func UnmarshalStruct(from, to interface{}) error {
	b, err := Marshal(from)
	if err != nil {
		return err
	}
	return Unmarshal(b, to)
}
