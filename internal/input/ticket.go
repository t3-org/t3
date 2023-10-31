package input

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateTicket struct {
	Fingerprint string `json:"fingerprint"`
	IsFiring    bool   `json:"is_firing"`
	StartedAt   int64  `json:"started_at"`
	EndedAt     *int64 `json:"ended_at"`

	IsSpam      bool    `json:"is_spam"`
	Level       *string `json:"level"`
	Description *string `json:"description"`
	SeenAt      *int64  `json:"seen_at"`
}

func (i *CreateTicket) Validate() error {
	return validation.ValidateStruct(&i) // TODO: write validations.

}

type UpdateTicket struct {
	Fingerprint *string `json:"fingerprint"`
	IsFiring    *bool   `json:"is_firing"`
	StartedAt   *int64  `json:"started_at"`
	EndedAt     *int64  `json:"ended_at"`

	IsSpam      *bool   `json:"is_spam"`
	Level       *string `json:"level"`
	Description *string `json:"description"`
	SeenAt      *int64  `json:"seen_at"`

	RemoveTags []string `json:"remove_tags"`
	AddTags    []string `json:"add_tags"`
}

func (i *UpdateTicket) Validate() error {
	return validation.ValidateStruct(&i) // TODO: update validations.

}

var _ validation.Validatable = &CreateTicket{}
var _ validation.Validatable = &UpdateTicket{}
