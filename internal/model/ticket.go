package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/kamva/gutil"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/helpers"
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
	Base              `json:",inline" yaml:",inline"`
	ID                string    `json:"id" yaml:"id"`
	GlobalFingerprint string    `json:"global_fingerprint" yaml:"global_fingerprint"`
	Fingerprint       string    `json:"fingerprint" yaml:"fingerprint"`
	Source            string    `json:"source" yaml:"source"` // Source of the alert. Maybe remove this field in the next versions.
	Raw               *string   `json:"raw" yaml:"raw"`       // the raw alert content. (optional)
	Annotations       StringMap `json:"annotations" yaml:"annotations"`
	IsFiring          bool      `json:"is_firing" yaml:"is_firing"`
	StartedAt         int64     `json:"started_at" yaml:"started_at"` // unix milliseconds.
	EndedAt           *int64    `json:"ended_at" yaml:"ended_at"`     // unix milliseconds.
	Values            StringMap `json:"values" yaml:"values"`
	GeneratorUrl      *string   `json:"generator_url"`
	IsSpam            bool      `json:"is_spam" yaml:"is_spam"`
	Severity          *string   `json:"severity" yaml:"severity"`
	Title             string    `json:"title" yaml:"title"`
	Description       *string   `json:"description" yaml:"description"`
	SeenAt            *int64    `json:"seen_at" yaml:"seen_at"` // unix milliseconds.

	// Internal labels start with "_". API can not touch(edit,remove...) internal labels.
	Labels map[string]string `sql:"" json:"labels" yaml:"labels"` // Set sql:"" to ignore this code generation scripts for DB.
}

func (m *Ticket) Create(in *input.CreateTicket) error {
	m.ID = genId()
	m.GlobalFingerprint = helpers.GlobalFingerprint(in.StartedAt, in.Fingerprint)
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
	if in.IsFromTicketGenerator {
		return m.patchTicketGeneratorUpdates(in)
	}
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
	// We let user update startedAt, but do not update globalFingerprint. because AlertManager never update startedAt field.
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

func (m *Ticket) patchTicketGeneratorUpdates(in *input.PatchTicket) error {
	// Update the firing state just if it's not resolved manually already.
	if m.IsFiring && in.IsFiring != nil { // Just update firing state if the ticket is firing.
		m.IsFiring = *in.IsFiring
	}

	if m.EndedAt == nil && in.EndedAt != nil {
		m.EndedAt = in.EndedAt
	}

	if in.IsFiring != nil && *in.IsFiring && in.Values != nil {
		m.Values = in.Values
	}

	m.Touch()
	return nil
}

func (m *Ticket) TitleMarkdown(isFiringOnTicketGenerator *bool) string {
	title := "__🔥🔥 Firing Ticket__" // when is firing

	if !m.IsFiring {
		title = "__🥳🥳 Resolved Ticket__" // When we resolve the ticket.
		if isFiringOnTicketGenerator != nil {
			title = "__🔥🔥 Firing state is continuing on the ticket generator__" // when is resolved by us, but firing on ticket generator
			if !*isFiringOnTicketGenerator {
				title = "__🥳🥳 Resolved on the ticket generator" // When resolve on the ticket genrator
			}
		}
	}
	return title
}

// Markdown returns the ticket data in markdown format.
// TODO: Use a user-defined template instead of this function.
func (m *Ticket) Markdown(verbose bool) string {
	b := strings.Builder{}
	w := func(val string, params ...any) {
		b.WriteString(fmt.Sprintf(val+"\n", params...))
	}

	isSpam := ""
	if m.IsSpam {
		isSpam = "(spam 🗑️)"
	}
	w("__%s__ %s", m.Title, isSpam)

	if verbose {
		firingEmoji := "🥳"
		if m.IsFiring {
			firingEmoji = "🔥"
		}

		w("- id: `%s`", m.ID)
		w("- is_firing: `%t`%s", m.IsFiring, firingEmoji)
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

	}

	w("- labels: \n")
	for k, v := range m.Labels {
		if !strings.HasPrefix(k, config.InternalLabelKeyPrefix) {
			w("	- `%s`: `%s`", k, v)
		}
	}

	w("- values: \n")
	for k, v := range m.Values {
		w("	- `%s`: `%s`", k, v)
	}

	if m.GeneratorUrl != nil {
		w("- [source link](%s)", *m.GeneratorUrl)
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

func TicketsMapByGlobalFingerprint(tickets ...*Ticket) map[string]*Ticket {
	res := make(map[string]*Ticket, len(tickets))
	for _, v := range tickets {
		res[v.GlobalFingerprint] = v
	}

	return res
}
