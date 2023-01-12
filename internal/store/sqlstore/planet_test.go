package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"space.org/space/internal/input"
	"space.org/space/internal/model"
	"space.org/space/internal/testbox"
)

func TestPlanetStore_Get(t *testing.T) {
	defer testbox.Global().TeardownIfPanic()

	s := store()
	ctx := context.Background()

	var p model.Planet
	require.NoError(t, p.Create(&input.CreatePlanet{
		Name: "name-abc",
		Code: "code-abc",
	}))

	require.NoError(t, s.Planet().Create(ctx, &p))
	defer func() { require.NoError(t, s.Planet().Delete(ctx, &p)) }()

	var p2 model.Planet
	require.NoError(t, p2.Create(&input.CreatePlanet{
		Name: "name-abc2",
		Code: "code-abc2",
	}))

	require.NoError(t, s.Planet().Create(ctx, &p2))
	defer func() { require.NoError(t, s.Planet().Delete(ctx, &p2)) }()

	res, err := s.Planet().Get(ctx, p.ID)
	require.NoError(t, err)
	require.Equal(t, &p, res)

	res2, err := s.Planet().Get(ctx, p2.ID)
	require.NoError(t, err)
	require.Equal(t, &p2, res2)
}

func TestPlanetStore_Create(t *testing.T) {
	defer testbox.Global().TeardownIfPanic()

	s := store()
	ctx := context.Background()

	var p model.Planet
	require.NoError(t, p.Create(&input.CreatePlanet{
		Name: "name-abc",
		Code: "code-abc",
	}))

	require.NoError(t, s.Planet().Create(ctx, &p))
	defer func() { require.NoError(t, s.Planet().Delete(ctx, &p)) }()

	count, err := s.Planet().Count(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 1, count)

	var p2 model.Planet
	require.NoError(t, p2.Create(&input.CreatePlanet{
		Name: "name-abc2",
		Code: "code-abc2",
	}))

	require.NoError(t, s.Planet().Create(ctx, &p2))
	defer func() { require.NoError(t, s.Planet().Delete(ctx, &p2)) }()

	count, err = s.Planet().Count(ctx, "")
	require.NoError(t, err)
	require.Equal(t, 2, count)
}
