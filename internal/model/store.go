//go:generate mockgen -source=store.go -destination=mock/store_gen.go -package=mockmodel
package model

import (
	"context"

	"github.com/kamva/hexa"
)

//nolint:unused
var dbStore Store

func SetStore(store Store) {
	dbStore = store
}

//nolint:unused
func store() Store {
	return dbStore
}

type Store interface {
	Health

	// DBLayer returns the database store layer.
	DBLayer() Store

	// System
	// @subStore
	System() SystemStore

	// Planet
	// @subStore
	Planet() PlanetStore

	// Place other store providers here.
}

type SystemStore interface {
	GetByName(ctx context.Context, name string) (*System, error)
	Save(ctx context.Context, system *System) error
	Delete(ctx context.Context, name string) error
}

type PlanetStore interface {
	Get(ctx context.Context, id string) (*Planet, error)
	GetByCode(ctx context.Context, code string) (*Planet, error)
	Create(ctx context.Context, m *Planet) error
	Update(ctx context.Context, m *Planet) error
	Delete(ctx context.Context, m *Planet) error
	Count(ctx context.Context, query string) (int, error)
	Query(ctx context.Context, query string, offset, limit uint64) ([]*Planet, error)
}

type Health interface {
	// HealthIdentifier
	// @noTracing
	HealthIdentifier() string
	LivenessStatus(ctx context.Context) hexa.LivenessStatus
	ReadinessStatus(ctx context.Context) hexa.ReadinessStatus
	HealthStatus(ctx context.Context) hexa.HealthStatus
}

var _ hexa.Health = Health(nil) // Assertion
