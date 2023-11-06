package model

import (
	"fmt"
	"strings"
	"time"

	"space.org/space/internal/input"
)

const (
	TicketLevelLow    = "low"
	TicketLevelMedium = "medium"
	TicketLevelHigh   = "high"
)

type Ticket struct {
	Base        `json:",inline" yaml:",inline"`
	ID          int64     `json:"id" yaml:"id"`
	Annotations StringMap `json:"annotations" yaml:"annotations"`
	Fingerprint string    `json:"fingerprint" yaml:"fingerprint"`
	IsFiring    bool      `json:"is_firing" yaml:"is_firing"`
	StartedAt   int64     `json:"started_at" yaml:"started_at"` // unix milliseconds.
	// TODO: Set to 0 if it's not set(for both of ended_at and seen_at fields).
	EndedAt *int64 `json:"ended_at" yaml:"ended_at"` // unix milliseconds.

	IsSpam      bool    `json:"is_spam" yaml:"is_spam"`
	Level       *string `json:"level" yaml:"level"`
	Description *string `json:"description" yaml:"description"`
	SeenAt      *int64  `json:"seen_at" yaml:"seen_at"` // unix milliseconds.

	// Internal labels start with "_". API can not touch(edit,remove...) internal labels.
	Labels map[string]string `sql:"" json:"labels" yaml:"labels"` // Set sql:"" to ignore this code generation scripts for DB.
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

func (m *Ticket) Create(in *input.CreateTicket) error {
	m.ID = genId()
	m.Fingerprint = in.Fingerprint
	m.IsFiring = in.IsFiring
	m.StartedAt = in.StartedAt
	m.EndedAt = in.EndedAt
	m.IsSpam = in.IsSpam
	m.Level = in.Level
	m.Description = in.Description
	m.SeenAt = in.SeenAt
	m.Labels = in.Labels
	m.Touch()
	return nil
}

func (m *Ticket) Patch(in *input.PatchTicket) error {
	if in.Fingerprint != nil {
		m.Fingerprint = *in.Fingerprint
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

	if in.Labels != nil {
		m.Labels = in.Labels
	}

	m.Touch()
	return nil
}
