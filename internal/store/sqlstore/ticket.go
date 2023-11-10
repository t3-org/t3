package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/helpers"
	"t3.org/t3/internal/model"
	"t3.org/t3/pkg/sqld"
)

const tableTicket = "tickets"

type ticketStore struct {
	s        SqlStore
	tbl      *sqld.Table
	labelTbl *sqld.Table
}

// newTicketStore returns new instance of the repository
func newTicketStore(store SqlStore) model.TicketStore {
	return &ticketStore{
		s:        store,
		tbl:      sqld.NewTable(tableTicket, ticketColumns, store.QueryBuilder),
		labelTbl: sqld.NewTable(tableTicketLabels, ticketLabelColumns, store.QueryBuilder),
	}
}

func (s *ticketStore) Get(ctx context.Context, id string) (*model.Ticket, error) {
	var ticket model.Ticket
	if err := s.tbl.FindByID(ctx, id, ticketFields(&ticket)); err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketNotFound))
	}
	if err := s.fetchLabels(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}
	return &ticket, nil
}

func (s *ticketStore) GetAllByFingerprint(ctx context.Context, fingerprints []string) ([]*model.Ticket, error) {
	if len(fingerprints) == 0 {
		return nil, nil
	}

	rows, err := s.tbl.Select(ctx).Where(sq.Eq{"fingerprint": fingerprints}).QueryContext(ctx)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, ticketFields)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return l, s.fetchLabelsForAll(ctx, l...)
}

func (s *ticketStore) FirstByTicketLabel(ctx context.Context, key, val string) (*model.Ticket, error) {
	var ticket model.Ticket
	err := s.tbl.First(ctx, ticketFields(&ticket),
		"id in (select ticket_id from ticket_labels where key=? and val=?)", key, val,
	)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	if err := s.fetchLabels(ctx, &ticket); err != nil {
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

	if err := s.fetchLabels(ctx, &ticket); err != nil {
		return nil, tracer.Trace(err)
	}

	return &ticket, nil
}

func (s *ticketStore) Create(ctx context.Context, ticket *model.Ticket) error {
	_, err := s.tbl.Create(ctx, ticketFields(ticket))
	if err != nil {
		return tracer.Trace(err)
	}

	return s.syncLabels(ctx, ticket)
}

func (s *ticketStore) Update(ctx context.Context, ticket *model.Ticket) error {
	_, err := s.tbl.Update(ctx, ticket.ID, ticketFields(ticket))
	if err != nil {
		return tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketNotFound))
	}

	return s.syncLabels(ctx, ticket)
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

	return l, s.fetchLabelsForAll(ctx, l...)
}

func (s *ticketStore) Delete(ctx context.Context, m *model.Ticket) error {
	_, err := s.tbl.Delete(ctx, m.ID)
	return tracer.Trace(err)
}

func (s *ticketStore) syncLabels(ctx context.Context, ticket *model.Ticket) error {
	// Remove not-existed labels.
	_, err := s.labelTbl.DeleteBuilder(ctx).
		Where(sq.Eq{"ticket_id": ticket.ID}).
		Where(sq.NotEq{"key": helpers.MapKeys(ticket.Labels)}).
		ExecContext(ctx)

	if err != nil {
		return tracer.Trace(err)
	}

	// Insert all labels
	if len(ticket.Labels) == 0 {
		return nil
	}

	b := s.s.QueryBuilder(ctx).Insert(tableTicketLabels).Columns(ticketLabelColumns...)
	b = sqld.SetValues(b, ticketLabelFields, model.LabelsFromMap(ticket.ID, ticket.Labels))

	// Add "on conflict(...) do update set field=excluded.field"
	expr, err := sqld.UpsertSuffix(
		sqld.OnConflictSet("ticket_id", "key"),
		sqld.Clauses(ticketLabelColumns, sqld.ExcludedColumns(ticketLabelColumns))...,
	)
	if err != nil {
		return tracer.Trace(err)
	}
	_, err = b.SuffixExpr(expr).ExecContext(ctx)
	return tracer.Trace(err)
}

func (s *ticketStore) fetchLabels(ctx context.Context, ticket *model.Ticket) error {
	rows, err := s.labelTbl.Select(ctx).Where(sq.Eq{"ticket_id": ticket.ID}).QueryContext(ctx)
	if err != nil {
		return tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, ticketLabelFields)
	if err != nil {
		return tracer.Trace(err)
	}
	ticket.Labels = model.LabelsMap(l...)
	return nil
}

func (s *ticketStore) fetchLabelsForAll(ctx context.Context, tickets ...*model.Ticket) error {
	if len(tickets) == 0 {
		return nil
	}

	ids := make([]string, len(tickets))
	ticketsMap := make(map[string]*model.Ticket, len(tickets))
	for i, v := range tickets {
		ids[i] = v.ID
		ticketsMap[v.ID] = v
	}

	rows, err := s.labelTbl.Select(ctx).Where(sq.Eq{"ticket_id": ids}).QueryContext(ctx)
	if err != nil {
		return tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, ticketLabelFields)
	if err != nil {
		return tracer.Trace(err)
	}

	for _, label := range l {
		t := ticketsMap[label.TicketID]
		if t.Labels == nil {
			t.Labels = make(map[string]string)
		}
		t.Labels[label.Key] = label.Val
	}
	return nil
}

var _ model.TicketStore = &ticketStore{}
