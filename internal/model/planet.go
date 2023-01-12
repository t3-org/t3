package model

import (
	"space.org/space/internal/input"
)

type Planet struct {
	Base `json:",inline"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"` // unique
}

func (m *Planet) Create(in *input.CreatePlanet) error {
	m.ID = genId()
	m.Name = in.Name
	m.Code = in.Code

	m.Touch()
	return nil
}

func (m *Planet) Update(in *input.UpdatePlanet) error {
	m.Name = in.Name
	m.Code = in.Code

	m.Touch()
	return nil
}
