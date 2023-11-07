package vcomplement

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// Translator is the ozzo-validation error translator.
	Translator interface {
		// Translate translate validation error.
		Translate(err error) (*TranslateBag, error)

		// Wrap translated messages in a hexa Error, if err is nil, return nil.
		WrapTranslationByError(err error) hexa.Error
	}

	// TranslateBag is the bag contains translated validation errors or error.
	TranslateBag struct {
		singleMessage string
		groupMessages map[string]*TranslateBag
	}
)

type hexaTranslator struct {
	t hexa.Translator
}

// NewHexaDriverErrorTranslator returns new instance of hexaTranslator
//that translate ozzo-validation errors.
func NewHexaDriverErrorTranslator(t hexa.Translator) Translator {
	return &hexaTranslator{t: t}
}

func (t *hexaTranslator) translateErr(err validation.Error) (string, error) {
	val, e := t.t.TranslateDefault(err.Code(), err.Message(), gutil.MapToKeyValue(err.Params())...)
	return val, tracer.Trace(e)
}

func (t *hexaTranslator) Translate(err error) (*TranslateBag, error) {
	bag := new(TranslateBag)

	if e, ok := gutil.CauseErr(err).(validation.Error); ok {
		msg, err := t.translateErr(e)
		bag.SetSingleMsg(msg)
		return bag, tracer.Trace(err)
	}

	if es, ok := gutil.CauseErr(err).(validation.Errors); ok {
		for k, e := range es {
			errBag, err := t.Translate(e)

			if err != nil {
				return nil, tracer.Trace(err)
			}

			bag.AddMsgToGroup(k, errBag)
		}

		return bag, nil
	}

	// otherwise return just empty bag and error (if error is internal
	// error, so user can detect it, otherwise get empty bag that to
	// detect that validation does not have any error).
	return bag, err
}

// Wrap bag translated error messages to hexa Error.
// if bag be empty, returns nil.
func (t *hexaTranslator) WrapTranslationByError(err error) hexa.Error {
	bag, err := t.Translate(err)

	if err != nil {
		return ErrInternalValidation.SetError(tracer.Trace(err))
	}

	// Bag can be nil (in case of valid data)
	if bag.IsEmpty() {
		return nil
	}

	return ErrValidationError.SetData(hexa.Map{"errors": bag.Map(true).(map[string]interface{})})
}

// SetSingleMsg set single error message.
func (t *TranslateBag) SetSingleMsg(msg string) {
	t.singleMessage = msg
}

// AddMsgToGroup add message to the group
func (t *TranslateBag) AddMsgToGroup(key string, msg *TranslateBag) {
	if t.groupMessages == nil {
		t.groupMessages = map[string]*TranslateBag{}
	}

	t.groupMessages[key] = msg
}

// IsEmpty specify that error bag is empty or not.
func (t *TranslateBag) IsEmpty() bool {
	return t.singleMessage == "" && len(t.groupMessages) == 0
}

// TOMap fucntion convert TranslateBag to a map[string]interface{}.
// but if bag just has a single message it check if forceMap is true,
// return map["error"]=<message>, otherwise returns just string message.
//
// possible values:
// - map[string]interface : when TranslateBag contains group of messages or
//	 contains single message with forceMap=true.
// - string: when TranslateBag contains singleMessage with forceMap=false.
// - nil : When TranslateBag does not hav single message nor group of messages.
func (t *TranslateBag) Map(forceMap bool) interface{} {
	if t.singleMessage != "" {
		if forceMap {
			return map[string]interface{}{"error": t.singleMessage}
		}

		return t.singleMessage
	}

	messages := map[string]interface{}{}

	for k, v := range t.groupMessages {
		messages[k] = v.Map(false)
	}

	return messages
}

// TValidate get a translator and validatable interface, validate and return hexa error.
func TValidateErr(t Translator, err error) error {
	return t.WrapTranslationByError(err)
}

// TValidate get a translator and validatable interface, validate and return hexa error.
func TValidate(t Translator, v validation.Validatable) error {
	return TValidateErr(t, v.Validate())
}

// TValidateByHexa validate by provided translator and check to detect right driver.
func TValidateByHexa(t hexa.Translator, v validation.Validatable) error {
	return TValidate(NewHexaDriverErrorTranslator(t.(hexa.Translator)), v)
}

// TranslateByHexa validate by provided translator and check to detect right driver.
func TranslateByHexa(t hexa.Translator, err error) error {
	return TValidateErr(NewHexaDriverErrorTranslator(t.(hexa.Translator)), err)
}

// Assert hexaTranslator implements the Translator.
var _ Translator = &hexaTranslator{}
