package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//--------------------------------
// Iran National ID
//--------------------------------

// ErrIRNationalIDInvalid is the default  IRNationalID validation rules error.
var ErrIRNationalIDInvalid = validation.NewError("validation_ir_national_id_invalid", "Iran national ID is invalid")

// IRNationalID is the Iran national ID validator.
var IRNationalID = validation.By(func(value interface{}) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	if str, ok := value.(string); ok && isValidIranianNationalCode(str) {
		return nil
	}

	return ErrIRNationalIDInvalid
})

func isValidIranianNationalCode(input string) bool {
	for i := 0; i < 10; i++ {
		if input[i] < '0' || input[i] > '9' {
			return false
		}
	}
	check := int(input[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(input[i]-'0') * (10 - i)
	}
	sum %= 11
	return (sum < 2 && check == sum) || (sum >= 2 && check+sum == 11)
}
