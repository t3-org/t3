package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//--------------------------------
// Iran Phone number.
//--------------------------------

// ErrInvalidID is the default invalid id validation error
var ErrInvalidID = validation.NewError("validation_invalid_id", "ID value is invalid")

type IDValidator func(val interface{}) error

// ID function return new validator to validate id.
func ID(idValidator IDValidator) validation.Rule {
	return validation.By(func(value interface{}) error {
		value, isNil := validation.Indirect(value)
		if isNil || validation.IsEmpty(value) {
			return nil
		}

		if err := idValidator(value); err != nil {
			return ErrInvalidID.SetMessage(err.Error())
		}

		return nil
	})
}
