package api

import (
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"space.org/space/internal/app"
	"space.org/space/internal/input"
)

func (api *API) registerPlanetRoutes(planets *echo.Group, res *planetResource) {
	planets.GET("/code/:code", res.GetByCode).Name = "planets::getByCode"
	hecho.ResourceAPI(planets, res, "planets")
}

type planetResource struct {
	Resource
	app app.App
}

func (r *planetResource) GetByCode(c echo.Context) error {
	dto, err := r.app.GetPlanetByCode(c.Request().Context(), c.Param("code"))
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessGetPlanet.SetData(dto))
}

func (r *planetResource) Create(c echo.Context) error {
	var in input.CreatePlanet
	if err := c.Bind(&in); err != nil {
		return tracer.Trace(err)
	}

	dto, err := r.app.CreatePlanet(c.Request().Context(), &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessCreatePlanet.SetData(dto))
}

func (r *planetResource) Update(c echo.Context) error {
	var in input.UpdatePlanet
	if err := c.Bind(&in); err != nil {
		return tracer.Trace(err)
	}

	dto, err := r.app.UpdatePlanet(c.Request().Context(), r.ID(c), &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessUpdatePlanet.SetData(dto))
}

func (r *planetResource) Delete(c echo.Context) error {
	err := r.app.DeletePlanet(c.Request().Context(), r.ID(c))
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessDeletePlanet)
}

func (r *planetResource) Query(c echo.Context) error {
	_, page, perPage := r.QueryAndPaginationParams(c)
	res, err := r.app.QueryPlanets(c.Request().Context(), c.QueryParam("query"), page, perPage)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessQueryPlanet.SetData(res))
}

var _ hecho.CreateResource = &planetResource{}
var _ hecho.UpdateResource = &planetResource{}
var _ hecho.DeleteResource = &planetResource{}
var _ hecho.QueryResource = &planetResource{}
