package input

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"t3.org/t3/internal/config"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/helpers"
)

type CreateTicket PatchTicket

func (i *CreateTicket) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.Fingerprint, validation.Required),
		validation.Field(&i.Source, validation.Required),
		validation.Field(&i.IsFiring, validation.NotNil),
		validation.Field(&i.StartedAt, validation.Required, validation.Min(1)),
		validation.Field(&i.EndedAt, validation.Min(1)),
		validation.Field(&i.IsSpam, validation.NotNil),
		validation.Field(&i.Severity, validation.NotNil, validation.In("low", "medium", "high")),
		validation.Field(&i.Title, validation.Required),
		validation.Field(&i.SeenAt, validation.Min(1)),
	)
}

type PatchTicket struct {
	isCreation bool
	// Set true when this input it's from a ticket generator(grafana alertManager,...)
	// For manual created tickets we don't have any ticketGenerator and don't need to set
	// this field to true.
	IsFromTicketGenerator bool `json:"-" yaml:"-"`

	// In patch requests, we'll ignore the fingerprint field. in creation
	// requests it's required.
	Fingerprint string `json:"fingerprint" yaml:"fingerprint"`

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
}

func (i *PatchTicket) SetIsCreation(isCreation bool) {
	i.isCreation = isCreation
	if i.isCreation {
		i.SetCreationDefaults()
	}
}

func (i *PatchTicket) SetCreationDefaults() {
	if i.IsSpam == nil {
		i.IsSpam = gutil.NewBool(false)
	}

	if i.Severity == nil {
		i.Severity = gutil.NewString("low")
	}
}

func (i *PatchTicket) Validate() error {
	if i.isCreation { // Another validation for creations.
		val := CreateTicket(*i)
		return val.Validate()
	}

	// Patch validation
	return validation.ValidateStruct(i,
		validation.Field(&i.StartedAt, validation.Min(1)),
		validation.Field(&i.EndedAt, validation.Min(1)),
		validation.Field(&i.Severity, validation.In("low", "medium", "high")),
		validation.Field(&i.SeenAt, validation.Min(1)),
	)
}

func RemoveInternalLabels(values map[string]string) {
	for k, _ := range values {
		if len(k) >= len(config.InternalLabelKeyPrefix) && strings.Index(k, config.InternalLabelKeyPrefix) == 0 {
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

func (i *BatchUpsertTickets) GlobalFingerprints() ([]string, error) {
	res := make([]string, len(i.Tickets))
	for i, v := range i.Tickets {
		f := helpers.GlobalFingerprint(v.StartedAt, v.Fingerprint)
		if f == "" {
			return nil, apperr.ErrTicketRequiredFieldsMissing.SetReportData(hexa.Map{
				"required_field": v.StartedAt,
				"fingerprint":    v.Fingerprint,
			})
		}
		res[i] = f
	}
	return res, nil
}

func (i *BatchUpsertTickets) RemoveInternalLabels() {
	for _, v := range i.Tickets {
		RemoveInternalLabels(v.Labels)
	}
}

var _ validation.Validatable = &CreateTicket{}
var _ validation.Validatable = &PatchTicket{}
var _ validation.Validatable = &BatchUpsertTickets{}
