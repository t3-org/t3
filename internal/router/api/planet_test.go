package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"space.org/space/internal/model"
	"space.org/space/internal/testbox"
)

var (
	planet = model.Planet{
		ID:   "1",
		Name: "abc-name",
		Code: "abc",
		Base: model.Base{
			CreatedAt: 1665660993990,
			UpdatedAt: 1665660993990,
		},
	}

	planetjson = `{"created_at":1665660993990,"updated_at":1665660993990,"id":"1","name":"abc-name","code":"abc"}`
)

func TestPlanetResource_GetByCode(t *testing.T) {
	_, a := setup(t)
	defer testbox.Global().TeardownIfPanic()

	a.EXPECT().GetPlanetByCode(_any, gomock.Eq("abc")).Return(&planet, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetPath("/planets/code/:code")
	c.SetParamNames("code")
	c.SetParamValues("abc")
	require.NoError(t, newPlanet(a).GetByCode(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, replyjson(RespSuccessGetPlanet.ID(), planetjson), rec.Body.String())
}
