package input

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreatePlanet struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (i *CreatePlanet) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Max(255)),
		validation.Field(&i.Code, validation.Required, is.DNSName, validation.Max(255)),
	)
}

type UpdatePlanet struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (i *UpdatePlanet) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required, validation.Max(255)),
		validation.Field(&i.Code, validation.Required, is.DNSName, validation.Max(255)),
	)
}

var _ validation.Validatable = &CreatePlanet{}
var _ validation.Validatable = &UpdatePlanet{}
