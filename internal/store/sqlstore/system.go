package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/pkg/sqld"
)

const tableSystems = "systems"

type systemStore struct {
	sqlStore SqlStore
	tbl      *sqld.Table
}

// newSystemStore returns new instance of the systemStore
func newSystemStore(store SqlStore) model.SystemStore {
	return &systemStore{
		sqlStore: store,
		tbl:      sqld.NewTable(tableSystems, systemColumns, store.QueryBuilder),
	}
}

func (s *systemStore) GetByName(ctx context.Context, name string) (*model.System, error) {
	var sys model.System
	if err := s.tbl.First(ctx, systemFields(&sys), sq.Eq{"name": name}); err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrSystemPropertyNotFound.SetData(hexa.Map{"name": name})))
	}
	return &sys, nil
}

func (s *systemStore) Save(ctx context.Context, system *model.System) error {
	_, err := s.tbl.Upsert(ctx, systemFields(system), sqld.OnConflictSET("name"))
	return tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrSystemPropertyNotFound))
}

func (s *systemStore) Update(ctx context.Context, system *model.System) error {
	_, err := s.tbl.UpdateBuilder(ctx, systemFields(system)).
		Where(sq.Eq{"name": system.Name}).
		ExecContext(ctx)

	return tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrSystemPropertyNotFound))
}

var _ model.SystemStore = &systemStore{}
