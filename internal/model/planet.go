package model

import (
	"github.com/kamva/gutil"
	"space.org/space/internal/input"
)

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