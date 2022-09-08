package sqlstore

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/kamva/tracer"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/pkg/sqld"
)

const tablePlanet = "planets"

type planetStore struct {
	s   SqlStore
	tbl *sqld.Table
}

// newPlanetStore returns new instance of the repository
func newPlanetStore(store SqlStore) model.PlanetStore {
	return &planetStore{
		s:   store,
		tbl: sqld.NewTable(tablePlanet, planetColumns, store.QueryBuilder),
	}
}

func (s *planetStore) Get(ctx context.Context, id string) (*model.Planet, error) {
	var planet model.Planet
	if err := s.tbl.FindByID(ctx, id, planetFields(&planet)); err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrPlanetNotFound))
	}

	return &planet, nil
}

func (s *planetStore) GetByCode(ctx context.Context, code string) (*model.Planet, error) {
	var planet model.Planet
	err := s.tbl.First(ctx, planetFields(&planet), sq.Eq{"code": code})
	if err != nil {
		return nil, tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrPlanetNotFound))
	}

	return &planet, nil
}

func (s *planetStore) Create(ctx context.Context, planet *model.Planet) error {
	_, err := s.tbl.Create(ctx, planetFields(planet))
	return tracer.Trace(err)
}

func (s *planetStore) Update(ctx context.Context, planet *model.Planet) error {
	_, err := s.tbl.Update(ctx, planet.ID, planetFields(planet))
	if err != nil {
		return tracer.Trace(sqld.ReplaceNotFound(err, apperr.ErrPlanetNotFound))
	}

	return nil
}

func (s *planetStore) Count(ctx context.Context, query string) (int, error) {
	var count int
	err := s.tbl.Count(ctx).Scan(&count)
	return count, err
}

func (s *planetStore) Query(ctx context.Context, query string, offset, limit uint64) ([]*model.Planet, error) {
	rows, err := s.tbl.Select(ctx).Limit(limit).Offset(offset).QueryContext(ctx)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	l, err := sqld.ScanRows(rows, planetFields)
	if err != nil {
		return nil, tracer.Trace(err)
	}
	return l, nil
}

func (s *planetStore) Delete(ctx context.Context, m *model.Planet) error {
	_, err := s.tbl.Delete(ctx, m.ID)
	return tracer.Trace(err)
}

var _ model.PlanetStore = &planetStore{}
