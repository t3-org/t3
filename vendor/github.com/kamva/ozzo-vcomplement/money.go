package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

//--------------------------------
// Iran postal code
//--------------------------------

// ErrMoneyInvalid is the default  Money validation rules error.
var ErrMoneyInvalid = validation.NewError("validation_money_invalid", "Money value is invalid")

// Money is the money validation rule. values should be positive decimals values.
var Money = validation.Match(regexp.MustCompile("^(0|[1-9]\\d*)?(\\.\\d+)?$")).ErrorObject(ErrMoneyInvalid)
