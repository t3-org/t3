package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"reflect"
)

// ErrInvalidConfirmation is the error that returns in case of invalid confirmation value.
var ErrInvalidConfirmation = validation.NewError("validation_invalid_confirmation", "confirmation value is invalid")

func Confirm(value interface{}) ConfirmRule {
	return ConfirmRule{
		realValue: value,
		err:       ErrInvalidConfirmation,
	}
}

// ConfirmRule is a validation rule that validates if a value can confirm another value.
type ConfirmRule struct {
	realValue interface{}
	err       validation.Error
}

// Validate checks if the given value is valid or not.
func (r ConfirmRule) Validate(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	realVal,isNil:=validation.Indirect(r.realValue)


	if isNil || !reflect.DeepEqual(realVal, value)  {
		return r.err
	}

	return nil
}

// Error sets the error message for the rule.
func (r ConfirmRule) Error(message string) ConfirmRule {
	r.err = r.err.SetMessage(message)
	return r
}

// ErrorObject sets the error struct for the rule.
func (r ConfirmRule) ErrorObject(err validation.Error) ConfirmRule {
	r.err = err
	return r
}
