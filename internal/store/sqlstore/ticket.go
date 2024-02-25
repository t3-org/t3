package sqlstore

import (
	"context"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/sets"
	"t3.org/t3/internal/config"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/helpers"
	"t3.org/t3/internal/model"
	"t3.org/t3/pkg/sqld"
)

const tableTicket = "tickets"

var queryableTicketCols = sets.New("id", "fingerprint", "source", "is_firing", "started_at", "ended_at", "is_spam", "severity", "seen_at")
var queryableDateCols = sets.New("started_at", "ended_at", "seen_at")
var queryableBoolCols = sets.New("is_firing", "is_spam")

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
func (s *ticketStore) Count(ctx context.Context, selector labels.Selector) (int, error) {
	var count int
	q := s.tbl.Count(ctx)
	var err error
	builder, err := s.makeQuery(ctx, q, selector)
	if err != nil {
		return 0, tracer.Trace(err)
	}
	q = *builder

	err = q.Scan(&count)
	return count, err
}

//nolint:revive // To disable unused query param issue.
func (s *ticketStore) Query(ctx context.Context, selector labels.Selector, offset, limit uint64) ([]*model.Ticket, error) {
	// Base query
	q := s.tbl.Select(ctx).Limit(limit).Offset(offset).OrderBy("started_at DESC")

	// Apply query
	builder, err := s.makeQuery(ctx, q, selector)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	q = *builder

	// Fetch
	rows, err := q.QueryContext(ctx)

	l, err := sqld.ScanRows(rows, ticketFields)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return l, s.fetchLabelsForAll(ctx, l...)
}

func (s *ticketStore) makeQuery(ctx context.Context, b sq.SelectBuilder, selector labels.Selector) (*sq.SelectBuilder, error) {
	requirements, selectable := selector.Requirements()
	if !selectable || len(requirements) == 0 {
		return &b, nil
	}

	builder, err := s.makeTicketConditions(b, requirements)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return s.makeLabelConditions(ctx, *builder, requirements)
}

func (s *ticketStore) makeTicketConditions(b sq.SelectBuilder, requirements labels.Requirements) (*sq.SelectBuilder, error) {
	// Add ticket conditions
	conditions := sq.Or{}
	for _, r := range requirements {
		key := strings.TrimPrefix(r.Key(), config.QueryTicketFieldsPrefix)
		if strings.Index(r.Key(), config.QueryTicketFieldsPrefix) != 0 {
			continue
		}

		// Make sure the col is in the white list:
		if !queryableTicketCols.Has(key) {
			return nil, apperr.ErrInvalidQuery.SetData(hexa.Map{"key": r.Key()})
		}

		values := make([]any, len(r.Values().List()))

		// convert values types if needed.
		for i, val := range r.Values().List() {
			var res any = val
			var err error
			if queryableDateCols.Has(key) {
				res, err = time.Parse(time.DateTime, val)
				if err != nil {
					return nil, apperr.ErrInvalidQuery.SetError(err)
				}

			} else if queryableBoolCols.Has(key) {
				res, err = strconv.ParseBool(val)
				if err != nil {
					return nil, apperr.ErrInvalidQuery.SetError(err)
				}
			}

			values[i] = res
		}

		var op sq.Sqlizer
		switch r.Operator() {
		case selection.Equals, selection.DoubleEquals, selection.In:
			op = sq.Eq{key: values}
		case selection.NotEquals, selection.NotIn:
			op = sq.NotEq{key: r.Values().List()}
		case selection.Exists:
			op = sq.NotEq{key: nil}
		default:
			return nil, apperr.ErrInvalidQuery.SetData(hexa.Map{"key": r.Key(), "operator": r.Operator()})
		}
		conditions = append(conditions, op)
	}

	if len(conditions) == 0 {
		return &b, nil
	}

	b = b.Where(conditions)
	return &b, nil
}

func (s *ticketStore) makeLabelConditions(ctx context.Context, b sq.SelectBuilder, requirements labels.Requirements) (*sq.SelectBuilder, error) {
	// Add labels conditions
	keys := sets.New[string]()
	conditions := sq.Or{}
	for _, r := range requirements {
		key := r.Key()
		if strings.Index(key, config.QueryTicketFieldsPrefix) == 0 {
			continue
		}

		var op sq.Sqlizer
		switch r.Operator() {
		case selection.Equals, selection.DoubleEquals, selection.In:
			op = sq.Eq{"key": key, "val": r.Values().List()}
		case selection.NotEquals, selection.NotIn:
			op = sq.NotEq{"key": key, "val": r.Values().List()}
		case selection.Exists:
			op = sq.NotEq{"key": key}
		default:
			return nil, apperr.ErrInvalidQuery.SetData(hexa.Map{"key": r.Key(), "operator": r.Operator()})
		}

		conditions = append(conditions, op)
		keys.Insert(r.Key())
	}

	if len(conditions) == 0 {
		return &b, nil
	}

	// make the labels query
	labelsQuery := s.labelTbl.Select(ctx, "1").Prefix("EXISTS(").Suffix(")").
		Where("ticket_id = tickets.id").
		Where(conditions).
		GroupBy("ticket_id").
		Having(sq.Eq{"count(*)": len(keys)})

	b = b.Where(labelsQuery) // and labels query
	return &b, nil
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

type TicketQueryMaker struct {
	queryCols       []string
	ticketColPrefix string
}
