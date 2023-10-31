package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/pkg/sqld"
)

const tableTicket = "tickets"
const tableTicketTags = "ticket_tags"

type ticketStore struct {
	s       SqlStore
	tbl     *sqld.Table
	tagsTbl *sqld.Table
}

// newTicketStore returns new instance of the repository
func newTicketStore(store SqlStore) model.TicketStore {
	return &ticketStore{
		s:       store,
		tbl:     sqld.NewTable(tableTicket, ticketColumns, store.QueryBuilder),
		tagsTbl: sqld.NewTable(tableTicketTags, ticketTagColumns, store.QueryBuilder),
	}
}

func (s *ticketStore) Get(ctx context.Context, id int64) (*model.Ticket, error) {
	var ticket model.Ticket
	if err := s.tbl.FindByID(ctx, id, ticketFields(&ticket)); err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketNotFound))
	}
	if err := s.fetchTags(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}
	return &ticket, nil
}

func (s *ticketStore) GetByCode(ctx context.Context, code string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := s.tbl.First(ctx, ticketFields(&ticket), sq.Eq{"code": code})
	if err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketNotFound))
	}

	if err := s.fetchTags(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}

	return &ticket, nil
}

func (s *ticketStore) Create(ctx context.Context, ticket *model.Ticket) error {
	_, err := s.tbl.Create(ctx, ticketFields(ticket))
	if err != nil {
		return tracer.Trace(err)
	}

	return s.syncTags(ctx, ticket)
}

func (s *ticketStore) Update(ctx context.Context, ticket *model.Ticket) error {
	_, err := s.tbl.Update(ctx, ticket.ID, ticketFields(ticket))
	if err != nil {
		return tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketNotFound))
	}

	return s.syncTags(ctx, ticket)
}

//nolint:revive // To disable unused query param issue.
func (s *ticketStore) Count(ctx context.Context, query string) (int, error) {
	var count int
	err := s.tbl.Count(ctx).Scan(&count)
	return count, err
}

//nolint:revive // To disable unused query param issue.
func (s *ticketStore) Query(ctx context.Context, query string, offset, limit uint64) ([]*model.Ticket, error) {
	rows, err := s.tbl.Select(ctx).Limit(limit).Offset(offset).QueryContext(ctx)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, ticketFields)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return l, nil
}

func (s *ticketStore) Delete(ctx context.Context, m *model.Ticket) error {
	_, err := s.tbl.Delete(ctx, m.ID)
	return tracer.Trace(err)
}

func (s *ticketStore) syncTags(ctx context.Context, ticket *model.Ticket) error {
	// Remove not-existed tags.
	_, err := s.tagsTbl.DeleteBuilder(ctx).
		Where(sq.Eq{"ticket_id": ticket.ID}).
		Where(sq.NotEq{"term": ticket.Tags}).
		ExecContext(ctx)

	if err != nil {
		return tracer.Trace(err)
	}

	// Insert all tags
	if len(ticket.Tags) == 0 {
		return nil
	}
	b := s.s.QueryBuilder(ctx).Insert(tableTicketTags).Columns(ticketTagColumns...)
	for _, term := range ticket.Tags {
		b.Values(ticket.ID, term)
	}
	b.SuffixExpr(sq.Expr(sqld.OnConflictDoNothing("ticket_id", "term"))) // ignore already existed tags
	_, err = b.ExecContext(ctx)
	return tracer.Trace(err)
}

func (s *ticketStore) fetchTags(ctx context.Context, ticket *model.Ticket) error {
	rows, err := s.tagsTbl.Select(ctx).Where(sq.Eq{"ticket_id": ticket.ID}).QueryContext(ctx)
	if err != nil {
		return tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, ticketTagFields)
	if err != nil {
		return tracer.Trace(err)
	}
	ticket.Tags = model.Tags(l...)
	return nil
}

var _ model.TicketStore = &ticketStore{}
