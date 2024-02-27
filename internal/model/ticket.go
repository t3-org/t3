package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kamva/gutil"
	"t3.org/t3/internal/input"
)

const (
	TicketSeverityLow    = "low"
	TicketSeverityMedium = "medium"
	TicketSeverityHigh   = "high"
)

const (
	SourceGrafana = "grafana"
)

type Ticket struct {
	Base         `json:",inline" yaml:",inline"`
	ID           string    `json:"id" yaml:"id"`
	Fingerprint  string    `json:"fingerprint" yaml:"fingerprint"`
	Source       string    `json:"source" yaml:"source"` // Source of the alert.
	Raw          *string   `json:"raw" yaml:"raw"`       // the raw alert content. (optional)
	Annotations  StringMap `json:"annotations" yaml:"annotations"`
	IsFiring     bool      `json:"is_firing" yaml:"is_firing"`
	StartedAt    int64     `json:"started_at" yaml:"started_at"` // unix milliseconds.
	EndedAt      *int64    `json:"ended_at" yaml:"ended_at"`     // unix milliseconds.
	Values       StringMap `json:"values" yaml:"values"`
	GeneratorUrl *string   `json:"generator_url"`
	IsSpam       bool      `json:"is_spam" yaml:"is_spam"`
	Severity     *string   `json:"severity" yaml:"severity"`
	Title        string    `json:"title" yaml:"title"`
	Description  *string   `json:"description" yaml:"description"`
	SeenAt       *int64    `json:"seen_at" yaml:"seen_at"` // unix milliseconds.

	// Internal labels start with "_". API can not touch(edit,remove...) internal labels.
	Labels map[string]string `sql:"" json:"labels" yaml:"labels"` // Set sql:"" to ignore this code generation scripts for DB.
}

func (m *Ticket) Create(in *input.CreateTicket) error {
	m.ID = genId()
	m.Fingerprint = in.Fingerprint
	m.Source = *in.Source
	if in.Raw != nil {
		m.Raw = in.Raw
	}
	m.Annotations = in.Annotations
	m.IsFiring = *in.IsFiring
	m.StartedAt = *in.StartedAt
	m.EndedAt = in.EndedAt
	m.Values = in.Values
	if in.GeneratorUrl != nil {
		m.GeneratorUrl = in.GeneratorUrl
	}
	m.IsSpam = *in.IsSpam
	m.Severity = in.Severity
	m.Title = *in.Title
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
	if in.Raw != nil {
		m.Raw = in.Raw
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

	if in.Severity != nil {
		m.Severity = in.Severity
	}

	if in.Title != nil {
		m.Title = *in.Title
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
			m.Labels = make(StringMap)
		}
		gutil.ExtendStrMap(m.Labels, in.Labels, true)
	}

	m.Touch()
	return nil
}

func (m *Ticket) Markdown() string {
	b := strings.Builder{}
	w := func(val string, params ...any) {
		b.WriteString(fmt.Sprintf(val+"\n", params...))
	}

	w("- id: `%s`", m.ID)
	w("- title: %s", m.Title)
	w("- is_firing: `%t`", m.IsFiring)
	w("- is_spam: `%t`", m.IsSpam)

	if m.Severity != nil {
		w("- severity: `%s`", *m.Severity)
	} else {
		w("- severity: `null`")
	}

	if m.Description != nil {
		w("- description: `%s`", *m.Description)
	} else {
		w("- description: `null`")
	}

	w("- started_at: `%s`", time.UnixMilli(m.StartedAt).Format(time.RFC1123))
	if m.EndedAt != nil {
		w("- ended_at: `%s`", time.UnixMilli(*m.EndedAt).Format(time.RFC1123))
	} else {
		w("- ended_at: `null`")
	}
	if m.SeenAt != nil {
		w("- seen_at: `%s`", time.UnixMilli(*m.SeenAt).Format(time.RFC1123))
	} else {
		w("- seen_at: `null`")
	}

	w("- labels: \n")
	for k, v := range m.Labels {
		w("	- `%s`: `%s`", k, v)
	}

	w("- values: \n")
	for k, v := range m.Values {
		w("	- `%s`: `%s`", k, v)
	}

	if m.GeneratorUrl != nil {
		w("- [generator](%s)", *m.GeneratorUrl)
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
