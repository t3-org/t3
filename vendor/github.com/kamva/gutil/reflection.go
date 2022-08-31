package gutil

import (
	"errors"
	"reflect"
)

// IsNil function check value is nil or no. To check real value of interface
// is nil or not, should using reflection, check this
// https://play.golang.org/p/Isoo0CcAvr. Firstly check
// `val==nil` because reflection can not get value of
// zero val.
func IsNil(val interface{}) (result bool) {

	if val == nil {
		return true
	}

	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return
}

// IndirectType returns indirect value's reflection type.
func IndirectType(val interface{}) (reflect.Type, error) {
	if IsNil(val) {
		return nil, errors.New("value can not be nil")
	}

	rType := reflect.TypeOf(val)

	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
	}

	return rType, nil
}

// Indirect returns the value that the given interface or pointer references to.
// A boolean value is also returned to indicate if the value is nil or not (only
// applicable to interface, pointer, map, and slice). If the value is neither an
// interface nor a pointer, it will be returned back.
// reference: https://github.com/go-ozzo/ozzo-validation/blob/master/util.go
func IndirectValue(value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	switch kind {
	case reflect.Invalid:
		return nil
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return IndirectValue(rv.Elem().Interface())
	case reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		if rv.IsNil() {
			return nil
		}
	}

	return value
}

// ValuePtr returns pointer to the provided value.
// Provided value can be either by reference or by value.
func ValuePtr(v interface{}) (interface{}, error) {
	t, err := IndirectType(v)
	if err != nil {
		return nil, err
	}

	return reflect.New(t).Interface(), nil
}

// MustValuePtr returns pointer to the provided value and panic if
// occurred error.
func MustValuePtr(v interface{}) interface{} {
	pv, err := ValuePtr(v)
	if err != nil {
		panic(err)
	}
	return pv
}

// StructTags return struct all fields tags.
func StructTags(val interface{}) ([]reflect.StructTag, error) {
	rType, err := IndirectType(val)

	if err != nil {
		return nil, err
	}

	if rType.Kind() != reflect.Struct {
		return nil, errors.New("value must be a struct or pointer to struct")
	}

	tags := make([]reflect.StructTag, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		tags = append(tags, rType.Field(i).Tag)
	}

	return tags, nil
}

// InterfaceToSlice converts interface to slices.
// If provided value is not list, it will panic.
func InterfaceToSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceToSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// InterfaceDefault returns value if it's not nil, otherwise
// returns the default value.
func InterfaceDefault(val, def interface{}) interface{} {
	if IsNil(val) {
		return def
	}
	return val
}
