package vcomplement

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//--------------------------------
// Iran phone number.
//--------------------------------

// ErrIRPhoneInvalid is the default IRPhone validation rules error.
var ErrIRPhoneInvalid = validation.NewError("validation_ir_phone_invalid", "Iran phone number is invalid")

var (
	IRPhoneRegex   = regexp.MustCompile("^(\\+98|0)9\\d{9}$")
	IRPhone98Regex = regexp.MustCompile("^\\+989\\d{9}$")
	IRPhone0Regex  = regexp.MustCompile("^09\\d{9}$")
)

// IRPhone is the iran phone number validation rule.
var (
	IRPhone   = validation.Match(IRPhoneRegex).ErrorObject(ErrIRPhoneInvalid)
	IRPhone98 = validation.Match(IRPhone98Regex).ErrorObject(ErrIRPhoneInvalid)
	IRPhone0  = validation.Match(IRPhone0Regex).ErrorObject(ErrIRPhoneInvalid)
)
