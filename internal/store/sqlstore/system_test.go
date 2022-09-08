package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
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
	require.NoError(t, s.System().Save(context.Background(), &sys))

	var sys1 model.System
	sys1.Create("abc-1", "123-1")
	require.NoError(t, s.System().Save(context.Background(), &sys1))

	res, err := s.System().GetByName(ctx, "abc")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc")
	require.Equal(t, res.Value, "123")
	require.NotZero(t, res.CreatedAt)
	require.NotZero(t, res.UpdatedAt)

	res, err = s.System().GetByName(ctx, "abc-1")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc-1")
	require.Equal(t, res.Value, "123-1")

	res.Value = "def-1"
	require.NoError(t, s.System().Save(ctx, res))
	res, err = s.System().GetByName(ctx, "abc-1")
	require.NoError(t, err)
	require.Equal(t, res.Name, "abc-1")
	require.Equal(t, res.Value, "def-1")
}
