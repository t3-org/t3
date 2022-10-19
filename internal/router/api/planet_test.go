package api

import (
	"fmt"
	"net"
	"testing"

	"github.com/kamva/hexa/hurl"
	"github.com/stretchr/testify/require"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
	"space.org/space/internal/testbox"
)

type Resp struct {
	Code string
	Data interface{}
}

func client(t *testing.T) *hurl.Client {
	svc := service()
	netAddr := svc.HttpServer().Echo.Listener.Addr().(*net.TCPAddr)
	addr := fmt.Sprintf("http://localhost:%d/api/v1", netAddr.Port)
	logmode := hurl.LogModeAll
	if !svc.Config().Debug {
		logmode = hurl.LogModeNone
	}
	cli, err := hurl.NewClient(addr, logmode)
	require.NoError(t, err)
	return cli
}

func TestPlanetResource_GetByCode(t *testing.T) {
	setup(t)
	defer testbox.Global().TeardownIfPanic()
	cli := client(t)
	r, err := cli.PostJSON("planets", &input.CreatePlanet{
		Name: "ABC",
		Code: "abc",
	})
	require.NoError(t, err)
	defer hurl.Drain(r)

	var planet dto.Planet
	require.NoError(t, hurl.ResponseErrOrDecodeJson(r, &Resp{Data: &planet}))

	require.Equal(t, "ABC", planet.Name)
	require.Equal(t, "abc", planet.Code)
	require.NotZero(t, planet.UpdatedAt)
	require.NotZero(t, planet.CreatedAt)

	r, err = cli.Get("planets/code/" + planet.Code)
	require.Nil(t, err)
	defer hurl.Drain(r)

	var planet2 dto.Planet
	require.NoError(t, hurl.ResponseErrOrDecodeJson(r, &Resp{Data: &planet2}))

	require.Equal(t, planet, planet2)
}
