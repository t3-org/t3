package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

//--------------------------------
// Iran telephone number.
//--------------------------------

// ErrIRTelInvalid is the default IRTelephone validation rules error.
var ErrIRTelInvalid = validation.NewError("validation_ir_tel_invalid", "Iran telephone number is invalid")

// IRTel is the iran telephone number validation rule.
var (
	IRTel   = validation.Match(regexp.MustCompile("^(\\+98|0)\\d{10}$")).ErrorObject(ErrIRTelInvalid)
	IRTel98 = validation.Match(regexp.MustCompile("^\\+98\\d{10}$")).ErrorObject(ErrIRTelInvalid)
	IRTel0  = validation.Match(regexp.MustCompile("^0\\d{10}$")).ErrorObject(ErrIRTelInvalid)
)
