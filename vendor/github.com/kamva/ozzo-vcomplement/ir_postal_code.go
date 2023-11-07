package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

//--------------------------------
// Iran postal code
//--------------------------------

// ErrIRPostalCodeInvalid is the default  IRPostalCode validation rules error.
var ErrIRPostalCodeInvalid = validation.NewError("validation_ir_postal_code_invalid", "Iran postal code number is invalid")

// IRPostalCode is the postal code of Iran rule
//var IRPostalCode = validation.Match(regexp.MustCompile(`^\b(?!(\d)\1{3})[13-9]{4}[1346-9][013-9]{5}\b$`)).ErrorObject(ErrIRPostalCodeInvalid)
var IRPostalCode = validation.Match(regexp.MustCompile("^\\d{10}$")).ErrorObject(ErrIRPostalCodeInvalid)
