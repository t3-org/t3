//go:generate mockgen -source=app_iface.go -destination=mock/app_gen.go -package=mockapp
package app

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/pagination"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
)

// App is core of the project
type App interface {
	Health
	PlanetService
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

type PlanetService interface {
	GetPlanet(ctx context.Context, id string) (*dto.Planet, error)
	GetPlanetByCode(ctx context.Context, code string) (*dto.Planet, error)
	CreatePlanet(ctx context.Context, in input.CreatePlanet) (*dto.Planet, error)
	UpdatePlanet(ctx context.Context, id string, in input.UpdatePlanet) (*dto.Planet, error)
	DeletePlanet(ctx context.Context, id string) error
	QueryPlanets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error)
}
