package model

import (
	"context"

	"github.com/kamva/hexa"
)

var dbStore Store

func SetDBStore(store Store) {
	dbStore = store
}

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

type Health interface {
	// HealthIdentifier
	// @noTracing
	HealthIdentifier() string
	LivenessStatus(ctx context.Context) hexa.LivenessStatus
	ReadinessStatus(ctx context.Context) hexa.ReadinessStatus
	HealthStatus(ctx context.Context) hexa.HealthStatus
}

var _ hexa.Health = Health(nil) // Assertion
