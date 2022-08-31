package hrpc

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func StringVal(v *string) *wrapperspb.StringValue {
	if v == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: *v}
}

func BoolVal(v *bool) *wrapperspb.BoolValue {
	if v == nil {
		return nil
	}
	return &wrapperspb.BoolValue{Value: *v}
}

func Int32Val(v *int32) *wrapperspb.Int32Value {
	if v == nil {
		return nil
	}
	return &wrapperspb.Int32Value{Value: *v}
}

func Int32ValFromInt(v *int) *wrapperspb.Int32Value {
	if v == nil {
		return nil
	}
	return &wrapperspb.Int32Value{Value: int32(*v)}
}

func Int64Val(v *int64) *wrapperspb.Int64Value {
	if v == nil {
		return nil
	}
	return &wrapperspb.Int64Value{Value: *v}
}

func FloatVal(v *float32) *wrapperspb.FloatValue {
	if v == nil {
		return nil
	}
	return &wrapperspb.FloatValue{Value: *v}
}

func DoubleVal(v *float64) *wrapperspb.DoubleValue {
	if v == nil {
		return nil
	}
	return &wrapperspb.DoubleValue{Value: *v}
}

func TimestampVal(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
