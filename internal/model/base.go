package model

import "time"

type Base struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Base) Touch() {
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}

	a.UpdatedAt = now
}
