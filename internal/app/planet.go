package app

import (
	"context"

	"github.com/kamva/hexa/pagination"
	"github.com/kamva/tracer"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
	"space.org/space/internal/model"
)

func (a *appCore) GetPlanet(ctx context.Context, id string) (*dto.Planet, error) {
	return a.store.Planet().Get(ctx, id)
}

func (a *appCore) GetPlanetByCode(ctx context.Context, code string) (*dto.Planet, error) {
	return a.store.Planet().GetByCode(ctx, code)
}

func (a *appCore) CreatePlanet(ctx context.Context, in input.CreatePlanet) (*dto.Planet, error) {
	var planet model.Planet
	if err := planet.Create(in); err != nil {
		return nil, tracer.Trace(err)
	}

	if err := a.store.Planet().Create(ctx, &planet); err != nil {
		return nil, tracer.Trace(err)
	}
	return &planet, nil
}

func (a *appCore) UpdatePlanet(ctx context.Context, id string, in input.UpdatePlanet) (*dto.Planet, error) {
	planet, err := a.store.Planet().Get(ctx, id)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	if err := planet.Update(in); err != nil {
		return nil, tracer.Trace(err)
	}
	if err := a.store.Planet().Update(ctx, planet); err != nil {
		return nil, tracer.Trace(err)
	}
	return planet, nil
}

func (a *appCore) DeletePlanet(ctx context.Context, id string) error {
	planet, err := a.store.Planet().Get(ctx, id)
	if err != nil {
		return tracer.Trace(err)
	}

	return tracer.Trace(a.store.Planet().Delete(ctx, planet))
}

//nolint:revive
func (a *appCore) QueryPlanets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error) {
	// TODO: implement me
	panic("implement me")
}

var _ PlanetService = &appCore{}
