package model

import "time"

type Base struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (a *Base) Touch() {
	now := time.Now()
	if a.CreatedAt == 0 {
		a.CreatedAt = now.UnixMilli()
	}

	a.UpdatedAt = now.UnixMilli()
}
