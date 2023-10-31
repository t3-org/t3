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
	TicketService
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

type TicketService interface {
	GetTicket(ctx context.Context, id int64) (*dto.Ticket, error)
	CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error)
	PatchTicket(ctx context.Context, id int64, in *input.UpdateTicket) (*dto.Ticket, error)
	DeleteTicket(ctx context.Context, id int64) error
	QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error)
}
