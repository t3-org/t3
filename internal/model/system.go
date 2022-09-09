package model

import (
	"context"
	"time"

	"github.com/kamva/tracer"
)

// System keys

const (
// SystemHelloWorld = "HELLO_WORLD"
)

type SystemStore interface {
	GetByName(ctx context.Context, name string) (*System, error)
	Save(ctx context.Context, system *System) error
}

type System struct {
	Base  `json:",inline"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (s *System) Create(name string, value string) {
	s.Name = name
	s.Value = value
	s.Touch()
}

func (s *System) SetName(name string) {
	s.Name = name
	s.Touch()
}

func (s *System) SetValue(val string) {
	s.Value = val
	s.Touch()
}

func (s *System) Time() (*time.Time, error) {
	t, err := time.Parse(time.RFC3339, s.Value)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return &t, nil
}

func (s *System) SetTime(t time.Time) {
	s.SetValue(t.Format(time.RFC3339))
}
