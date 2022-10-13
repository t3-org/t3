package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	apperr "space.org/space/internal/error"
	"space.org/space/internal/model"
	"space.org/space/internal/registry"
	"space.org/space/internal/testbox"
)

func store() model.Store {
	return testbox.Global().Registry().Service(registry.ServiceNameStore).(model.Store)
}

func TestSystemStore_Save(t *testing.T) {
	defer testbox.Global().TeardownIfPanic()
	s := store()
	ctx := context.Background()

	var sys model.System
	sys.Create("abc", "123")
	require.NoError(t, s.System().Save(ctx, &sys))
	defer func() { require.NoError(t, s.System().Delete(ctx, sys.Name)) }()

	var sys2 model.System
	sys2.Create("abc-2", "123-2")
	require.NoError(t, s.System().Save(ctx, &sys2))
	defer func() { require.NoError(t, s.System().Delete(ctx, sys.Name)) }()

	res, err := s.System().GetByName(ctx, "abc")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc")
	require.Equal(t, res.Value, "123")
	require.NotZero(t, res.CreatedAt)
	require.NotZero(t, res.UpdatedAt)

	res, err = s.System().GetByName(ctx, "abc-2")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc-2")
	require.Equal(t, res.Value, "123-2")

	res.Value = "def-1"
	require.NoError(t, s.System().Save(ctx, res))
	res, err = s.System().GetByName(ctx, "abc-2")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc-2")
	require.Equal(t, res.Value, "def-1")
}

func TestSystemStore_Delete(t *testing.T) {
	defer testbox.Global().TeardownIfPanic()

	s := store()
	ctx := context.Background()

	var sys model.System
	sys.Create("abc", "123")
	require.NoError(t, s.System().Save(ctx, &sys))

	var sys2 model.System
	sys2.Create("abc-2", "123-2")
	require.NoError(t, s.System().Save(ctx, &sys2))

	require.NoError(t, s.System().Delete(ctx, sys.Name))
	res, err := s.System().GetByName(ctx, "abc")
	require.ErrorIs(t, err, apperr.ErrSystemPropertyNotFound)
	require.Nil(t, res)

	res, err = s.System().GetByName(ctx, "abc-2")
	require.NoError(t, err)
	require.NotNil(t, res)

	require.NoError(t, s.System().Delete(ctx, sys2.Name))
	res, err = s.System().GetByName(ctx, "abc")
	require.ErrorIs(t, err, apperr.ErrSystemPropertyNotFound)
	require.Nil(t, res)
}
