package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/pkg/sqld"
)

const tableTicketValues = "ticket_values"

type ticketKVStore struct {
	s       SqlStore
	tbl     *sqld.Table
	tagsTbl *sqld.Table
}

// newTicketStore returns new instance of the repository
func newTicketKVStore(store SqlStore) model.TicketKVStore {
	return &ticketKVStore{
		s:   store,
		tbl: sqld.NewTable(tableTicketValues, ticketKVColumns, store.QueryBuilder),
	}
}

func (s *ticketKVStore) Set(ctx context.Context, ticketID int64, key string, val string) error {
	kv := &model.TicketKV{
		TicketID: ticketID,
		Key:      key,
		Value:    val,
	}
	_, err := s.tbl.Upsert(ctx, ticketKVFields(kv), sqld.OnConflictSet("ticket_id", "key"))
	return tracer.Trace(err)
}

func (s *ticketKVStore) Val(ctx context.Context, id int64, key string) (string, error) {
	var kv model.TicketKV
	err := s.tbl.First(ctx, ticketKVFields(&kv), sq.Eq{"ticket_id": id, "key": key})
	if err != nil {
		return "", tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrTicketKVNotFound))
	}

	return kv.Value, nil
}
