package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"space.org/space/internal/input"
	"space.org/space/internal/model"
	"space.org/space/internal/model/mock"
	"space.org/space/internal/testbox"
)

func TestAppCore_CreatePlanet(t *testing.T) {
	defer testbox.Global().TeardownIfPanic()
	_, s, a := setup(t)
	now := time.Now().UnixMilli()

	_, _ = s, a
	in := input.CreatePlanet{
		Name: "abc",
		Code: "abc",
	}

	var res *model.Planet
	s.Planet().(*mockmodel.MockPlanetStore).EXPECT().Create(_any, _any).DoAndReturn(
		func(_ context.Context, m *model.Planet) error {
			res = m
			return nil
		},
	)

	p, err := a.CreatePlanet(context.Background(), in)
	require.NoError(t, err)
	require.Equal(t, p, res)
	require.NotEmpty(t, p.ID)
	require.GreaterOrEqual(t, p.CreatedAt, now)
	require.GreaterOrEqual(t, p.UpdatedAt, now)
	require.Equal(t, p.Name, in.Name)
	require.Equal(t, p.Code, in.Code)
}
