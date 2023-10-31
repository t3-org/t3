package model

import (
	"github.com/kamva/gutil"
	"space.org/space/internal/input"
)

const (
	TicketLevelLow    = "low"
	TicketLevelMedium = "medium"
	TicketLevelHigh   = "high"
)

type Ticket struct {
	Base        `json:",inline"`
	ID          int64  `json:"id"`
	Fingerprint string `json:"fingerprint"`
	IsFiring    bool   `json:"is_firing"`
	StartedAt   int64  `json:"started_at"`
	EndedAt     *int64 `json:"ended_at"`

	IsSpam      bool    `json:"is_spam"`
	Level       *string `json:"level"`
	Description *string `json:"description"`
	SeenAt      *int64  `json:"seen_at"`

	Tags []string `sql:"" json:"tags"`
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

	m.Touch()
	return nil
}

func (m *Ticket) Patch(in *input.UpdateTicket) error {
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

	if len(in.RemoveTags) != 0 {
		m.Tags = append(gutil.RemoveFromStrings(m.Tags, in.RemoveTags...), in.AddTags...)
	}

	m.Touch()
	return nil
}
