package gutil

import "errors"

// Must check if error is not nil, so panic, otherwise return value.
func Must(val interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return val
}

// AnyErr returns first error that is not nil.
func AnyErr(errors ...error) error {
	for _, err := range errors {
		if !IsNil(err) {
			return err
		}
	}

	return nil
}

// CauseErr returns the root error by unwrapping the errors chain
func CauseErr(err error) error {
	for errors.Unwrap(err) != nil {
		err = errors.Unwrap(err)
	}
	return err
}
