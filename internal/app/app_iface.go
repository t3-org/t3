//go:generate mockgen -source=app_iface.go -destination=mock/app_gen.go -package=mockapp
package app

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/pagination"
	"t3.org/t3/internal/dto"
	"t3.org/t3/internal/input"
)

// App is core of the project
type App interface {
	Health
	Seeder
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

type Seeder interface {
	Seed(ctx context.Context, count int32) error
}

type TicketService interface {
	UpsertTickets(ctx context.Context, in *input.BatchUpsertTickets) ([]*dto.Ticket, error)
	EditTicketUrlByKey(ctx context.Context, key, val string) (string, error)
	GetTicketByKey(ctx context.Context, key, val string) (*dto.Ticket, error)
	GetTicket(ctx context.Context, id string) (*dto.Ticket, error)
	CreateTicket(ctx context.Context, in *input.CreateTicket) (*dto.Ticket, error)
	PatchTicket(ctx context.Context, id string, in *input.PatchTicket) (*dto.Ticket, error)
	PatchTicketByLabel(ctx context.Context, key, val string, in *input.PatchTicket) (*dto.Ticket, error)
	DeleteTicket(ctx context.Context, id string) error
	QueryTickets(ctx context.Context, query string, page, perPage int) (*pagination.Pages, error)
}
