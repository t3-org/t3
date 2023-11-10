package input

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Channel struct {
	Name string `json:"name"` // Name of the channel (to find the channel instance that we should use)
	ID   string `json:"id"`   // id of the channel (e.g., in matrix it's the roomID).
}

type CreateTicket PatchTicket

func (i *CreateTicket) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Severity, validation.In("low", "medium", "high")),
	) // TODO: write validations.

}

type PatchTicket struct {
	isCreation bool

	Fingerprint string `json:"fingerprint" yaml:"fingerprint"` // In patch requests, we'll ignore this field.

	Source          *string           `json:"source" yaml:"source"`
	Raw             *string           `json:"raw" yaml:"raw"`
	Annotations     map[string]string `json:"annotations" yaml:"annotations"`
	SyncAnnotations bool              `json:"sync_annotations" yaml:"sync_annotations"` // if true, set annotations, otherwise upsert the provided annotations.
	IsFiring        *bool             `json:"is_firing" yaml:"is_firing"`
	StartedAt       *int64            `json:"started_at" yaml:"started_at"`
	EndedAt         *int64            `json:"ended_at" yaml:"ended_at"`
	Values          map[string]string `json:"values" yaml:"values"`
	GeneratorUrl    *string           `json:"generator_url" yaml:"generator_url"`
	IsSpam          *bool             `json:"is_spam" yaml:"is_spam"`
	Severity        *string           `json:"severity" yaml:"severity"`
	Title           *string           `json:"title" yaml:"title"`
	Description     *string           `json:"description" yaml:"description"`
	SeenAt          *int64            `json:"seen_at" yaml:"seen_at"`
	Labels          map[string]string `json:"labels" yaml:"labels"`
	SyncLabels      bool              `json:"sync_labels" yaml:"sync_labels"`
	Channel         *Channel          `json:"channel" yaml:"channel"`
}

func (i *PatchTicket) SetIsCreation(isCreation bool) {
	i.isCreation = isCreation
}

func (i *PatchTicket) Validate() error {
	if i.isCreation { // Another validation for creations.
		val := CreateTicket(*i)
		return val.Validate()
	}

	return validation.ValidateStruct(i,
		validation.Field(&i.Severity, validation.In("low", "medium", "high")),
	) // TODO: update validations.

}

func RemoveInternalLabels(values map[string]string) {
	for k, _ := range values {
		if len(k) != 0 && k[0] == '_' {
			delete(values, k)
		}
	}
}

type BatchUpsertTickets struct {
	Tickets []*PatchTicket `json:"tickets"`
}

func (i *BatchUpsertTickets) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Tickets, validation.Each(validation.Required)))
}

func (i *BatchUpsertTickets) Fingerprints() []string {
	res := make([]string, len(i.Tickets))
	for i, v := range i.Tickets {
		res[i] = v.Fingerprint
	}
	return res
}

func (i *BatchUpsertTickets) RemoveInternalLabels() {
	for _, v := range i.Tickets {
		RemoveInternalLabels(v.Labels)
	}
}

var _ validation.Validatable = &CreateTicket{}
var _ validation.Validatable = &PatchTicket{}
var _ validation.Validatable = &BatchUpsertTickets{}
