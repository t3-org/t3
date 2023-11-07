package vcomplement

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

//--------------------------------
// Iran phone number.
//--------------------------------

// ErrIRCreditCardInvalid is the default IRCreditCard validation rules error.
var ErrIRCardNumberInvalid = validation.NewError("validation_ir_card_number_invalid", "Iran card number is invalid")
var ErrIRAccountNumberInvalid = validation.NewError("validation_ir_account_number_invalid", "Iran account number is invalid")
var ErrIRIBANNumberInvalid = validation.NewError("validation_ir_iban_invalid", "Iran IBAN is invalid")

var (
	IRCardNumber    = validation.Match(regexp.MustCompile("^\\d{16}$")).ErrorObject(ErrIRCardNumberInvalid)
	IRAccountNumber = validation.Match(regexp.MustCompile("^(\\d{5,19})$")).ErrorObject(ErrIRAccountNumberInvalid)
	IRIBAN          = validation.Match(regexp.MustCompile("^IR[0-9]{24}$")).ErrorObject(ErrIRIBANNumberInvalid)
)
