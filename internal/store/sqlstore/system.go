package sqlstore

import (
	"context"

	"github.com/kamva/mgm/v3"
	"space.org/space/internal/model"
)

type systemStore struct {
	sqlStore SqlStore
}

func (s *systemStore) GetByName(ctx context.Context, name string) (*model.System, error) {
	// TODO implement me
	panic("implement me")
}

func (s *systemStore) Save(ctx context.Context, system *model.System) error {
	// TODO implement me
	panic("implement me")
}

func (s *systemStore) SaveModel(ctx context.Context, system mgm.Model) error {
	// TODO implement me
	panic("implement me")
}

// newSystemStore returns new instance of the systemStore
func newSystemStore(store SqlStore) model.SystemStore {
	return &systemStore{sqlStore: store}
}

var _ model.SystemStore = &systemStore{}
