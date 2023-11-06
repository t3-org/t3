package input

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Webhook struct {
	Channel   string `json:"channel"`
	ChannelID string `json:"channel_id"`
}

type CreateTicket struct {
	Fingerprint string `json:"fingerprint"`
	IsFiring    bool   `json:"is_firing"`
	StartedAt   int64  `json:"started_at"`
	EndedAt     *int64 `json:"ended_at"`

	IsSpam      bool              `json:"is_spam"`
	Level       *string           `json:"level"`
	Description *string           `json:"description"`
	SeenAt      *int64            `json:"seen_at"`
	Labels      map[string]string `json:"labels"`
	Webhook     *Webhook          `json:"webhook"`
}

func (i *CreateTicket) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Level, validation.In("low", "medium", "high")),
	) // TODO: write validations.

}

type PatchTicket struct {
	Fingerprint *string `json:"fingerprint"`
	IsFiring    *bool   `json:"is_firing"`
	StartedAt   *int64  `json:"started_at"`
	EndedAt     *int64  `json:"ended_at"`

	IsSpam      *bool   `json:"is_spam"`
	Level       *string `json:"level"`
	Description *string `json:"description"`
	SeenAt      *int64  `json:"seen_at"`

	Labels map[string]string

	Webhook *Webhook `json:"webhook"`
}

func (i *PatchTicket) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Level, validation.In("low", "medium", "high")),
	) // TODO: update validations.

}

var _ validation.Validatable = &CreateTicket{}
var _ validation.Validatable = &PatchTicket{}

func RemoveInternalLabels(labels map[string]string) {
	for k, _ := range labels {
		if len(k) != 0 && k[0] == '_' {
			delete(labels, k)
		}
	}
}
