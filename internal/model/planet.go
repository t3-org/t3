package model

import (
	"context"

	"github.com/kamva/gutil"
	"space.org/space/internal/input"
)

type PlanetStore interface {
	Get(ctx context.Context, id string) (*Planet, error)
	GetByCode(ctx context.Context, code string) (*Planet, error)
	Create(ctx context.Context, m *Planet) error
	Update(ctx context.Context, m *Planet) error
	Delete(ctx context.Context, m *Planet) error
	Count(ctx context.Context, query string) (int, error)
	Query(ctx context.Context, query string, offset, limit uint64) ([]*Planet, error)
}

type Planet struct {
	Base `json:",inline"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // unique
}

func (s *Planet) Create(in input.CreatePlanet) error {
	s.ID = gutil.UUID()
	s.Name = in.Name
	s.Code = in.Code

	s.Touch()
	return nil
}

func (s *Planet) Update(in input.UpdatePlanet) error {
	s.Name = in.Name
	s.Code = in.Code

	s.Touch()
	return nil
}
