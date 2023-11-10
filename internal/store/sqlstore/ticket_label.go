package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/model"
	"t3.org/t3/pkg/sqld"
)

const tableTicketLabels = "ticket_labels"

type ticketLabelsStore struct {
	s   SqlStore
	tbl *sqld.Table
}

// newTicketLabelStore returns new instance of the repository
func newTicketLabelStore(store SqlStore) model.TicketLabelStore {
	return &ticketLabelsStore{
		s:   store,
		tbl: sqld.NewTable(tableTicketLabels, ticketLabelColumns, store.QueryBuilder),
	}
}

func (s *ticketLabelsStore) Set(ctx context.Context, ticketID string, key string, val string) error {
	label := &model.TicketLabel{
		TicketID: ticketID,
		Key:      key,
		Val:      val,
	}
	_, err := s.tbl.Upsert(ctx, ticketLabelFields(label), sqld.OnConflictSet("ticket_id", "key"))
	return tracer.Trace(err)
}

func (s *ticketLabelsStore) Val(ctx context.Context, id string, key string) (string, error) {
	var label model.TicketLabel
	err := s.tbl.First(ctx, ticketLabelFields(&label), sq.Eq{"ticket_id": id, "key": key})
	if err != nil {
		return "", tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketKVNotFound))
	}

	return label.Val, nil
}
