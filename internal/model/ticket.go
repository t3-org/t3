package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kamva/gutil"
	"space.org/space/internal/input"
)

const (
	TicketLevelLow    = "low"
	TicketLevelMedium = "medium"
	TicketLevelHigh   = "high"
)

const (
	SourceGrafana = "grafana"
)

type Ticket struct {
	Base         `json:",inline" yaml:",inline"`
	ID           int64     `json:"id" yaml:"id"`
	Source       string    `json:"source"`    // Source of the alert.
	RawAlert     *string   `json:"raw_alert"` // the raw alert content. (optional)
	Fingerprint  string    `json:"fingerprint" yaml:"fingerprint"`
	Annotations  StringMap `json:"annotations" yaml:"annotations"`
	IsFiring     bool      `json:"is_firing" yaml:"is_firing"`
	StartedAt    int64     `json:"started_at" yaml:"started_at"` // unix milliseconds.
	EndedAt      *int64    `json:"ended_at" yaml:"ended_at"`     // unix milliseconds.
	Values       StringMap `json:"values" yaml:"values"`
	GeneratorUrl *string   `json:"generator_url"`
	IsSpam       bool      `json:"is_spam" yaml:"is_spam"`
	Level        *string   `json:"level" yaml:"level"`
	Description  *string   `json:"description" yaml:"description"`
	SeenAt       *int64    `json:"seen_at" yaml:"seen_at"` // unix milliseconds.

	// Internal labels start with "_". API can not touch(edit,remove...) internal labels.
	Labels map[string]string `sql:"" json:"labels" yaml:"labels"` // Set sql:"" to ignore this code generation scripts for DB.
}

func (m *Ticket) Create(in *input.CreateTicket) error {
	m.ID = genId()
	m.Source = *in.Source
	if in.RawAlert != nil {
		m.RawAlert = in.RawAlert
	}
	m.Fingerprint = in.Fingerprint
	m.Annotations = in.Annotations
	m.IsFiring = *in.IsFiring
	m.StartedAt = *in.StartedAt
	m.EndedAt = in.EndedAt
	m.Values = in.Values
	if in.GeneratorUrl != nil {
		m.GeneratorUrl = in.GeneratorUrl
	}
	m.IsSpam = *in.IsSpam
	m.Level = in.Level
	m.Description = in.Description
	m.SeenAt = in.SeenAt
	m.Labels = in.Labels

	m.Touch()
	return nil
}

func (m *Ticket) Patch(in *input.PatchTicket) error {
	if in.Source != nil {
		m.Source = *in.Source
	}
	if in.RawAlert != nil {
		m.RawAlert = in.RawAlert
	}

	if in.SyncAnnotations {
		m.Annotations = in.Annotations
	} else {
		if m.Annotations == nil {
			m.Annotations = make(StringMap)
		}
		gutil.ExtendStrMap(m.Annotations, in.Annotations, true)
	}

	if in.IsFiring != nil {
		m.IsFiring = *in.IsFiring
	}
	if in.StartedAt != nil {
		m.StartedAt = *in.StartedAt
	}
	if in.EndedAt != nil {
		m.EndedAt = in.EndedAt
	}

	if in.Values != nil {
		m.Values = in.Values
	}

	if in.GeneratorUrl != nil {
		m.GeneratorUrl = in.GeneratorUrl
	}

	if in.IsSpam != nil {
		m.IsSpam = *in.IsSpam
	}

	if in.Level != nil {
		m.Level = in.Level
	}
	if in.Description != nil {
		m.Description = in.Description
	}
	if in.SeenAt != nil {
		m.SeenAt = in.SeenAt
	}

	if in.SyncLabels {
		m.Labels = InternalLabels(m.Labels) // we'll keep internal labels.
		gutil.ExtendStrMap(in.Labels, in.Labels, true)
	} else {
		if m.Labels == nil {
			m.Annotations = make(StringMap)
		}
		gutil.ExtendStrMap(m.Annotations, in.Annotations, true)
	}

	m.Touch()
	return nil
}

func (m *Ticket) Markdown() string {
	b := strings.Builder{}
	w := func(val string, params ...any) {
		b.WriteString(fmt.Sprintf(val+"\n", params...))
	}

	w("- id: `%d`", m.ID)
	w("- is_firing: `%t`", m.IsFiring)
	w("- is_spam: `%t`", m.IsSpam)

	if m.Level != nil {
		w("- level: `%s`", *m.Level)
	} else {
		w("- level: `null`")
	}

	if m.Description != nil {
		w("- description: `%s`", *m.Description)
	} else {
		w("- level: `null`")
	}

	w("- started_at: `%s`", time.UnixMilli(m.StartedAt).Format(time.RFC3339))
	if m.EndedAt != nil {
		w("- ended_at: `%s`", time.UnixMilli(*m.EndedAt).Format(time.RFC3339))
	} else {
		w("- ended_at: `null`")
	}
	if m.SeenAt != nil {
		w("- seen_at: `%s`", time.UnixMilli(*m.SeenAt).Format(time.RFC3339))
	} else {
		w("- seen-at: `null`")
	}

	w("- labels: \n")
	for k, v := range m.Labels {
		w("	- %s: %s", k, v)
	}

	return b.String()
}

func InternalLabels(m map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range m {
		if len(k) != 0 && k[0] == '_' {
			res[k] = v
		}
	}
	return res
}

func TicketsMapByFingerprint(tickets ...*Ticket) map[string]*Ticket {
	res := make(map[string]*Ticket, len(tickets))
	for _, v := range tickets {
		res[v.Fingerprint] = v
	}

	return res
}
