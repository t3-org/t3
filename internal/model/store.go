//go:generate mockgen -source=store.go -destination=mock/store_gen.go -package=mockmodel
package model

import (
	"context"

	"github.com/kamva/hexa"
	"space.org/space/pkg/sqld"
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

	Txs() *sqld.Txs

	TruncateAllTables(ctx context.Context) error

	// System
	// @subStore
	System() SystemStore

	// Ticket
	// @subStore
	Ticket() TicketStore

	// TicketLabel
	// @subStore
	TicketLabel() TicketLabelStore

	// Place other store providers here.

}

type SystemStore interface {
	GetByName(ctx context.Context, name string) (*System, error)
	Save(ctx context.Context, system *System) error
	Delete(ctx context.Context, name string) error
}

type TicketStore interface {
	Get(ctx context.Context, id int64) (*Ticket, error)
	GetAllByFingerprint(ctx context.Context, fingerprints []string) ([]*Ticket, error)
	FirstByTicketLabel(ctx context.Context, key, val string) (*Ticket, error)
	GetByCode(ctx context.Context, code string) (*Ticket, error)
	Create(ctx context.Context, m *Ticket) error
	Update(ctx context.Context, m *Ticket) error
	Delete(ctx context.Context, m *Ticket) error
	Count(ctx context.Context, query string) (int, error)
	Query(ctx context.Context, query string, offset, limit uint64) ([]*Ticket, error)
}

type TicketLabelStore interface {
	Set(ctx context.Context, ticketID int64, key string, val string) error
	Val(ctx context.Context, ticketID int64, key string) (string, error)
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
