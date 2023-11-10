package api

import (
	"strings"

	hecho "github.com/kamva/hexa-echo"
	"github.com/labstack/echo/v4"
	"t3.org/t3/internal/app"
)

func (api *API) registerLabRoutes(lab *echo.Group, res *labResource) {
	lab.GET("/routes", res.Routes, api.middlewares.Debug).Name = "lab::routes"
	lab.GET("/ping", res.Pong).Name = "lab::ping"
}

type labResource struct {
	Resource
	e   *echo.Echo
	app app.App
}

func (r *labResource) Routes(c echo.Context) error {
	data := r.e.Routes()
	if q := c.QueryParam("query"); q != "" {
		data = r.filter(data, q)
	}
	return hecho.Write(c, RespSuccessGetRoutes.SetData(map[string]interface{}{
		"routes": data,
	}))
}

func (r *labResource) filter(routes []*echo.Route, query string) []*echo.Route {
	filtered := make([]*echo.Route, 0)
	for _, route := range routes {
		if r.match(query, route.Name, route.Method, route.Path) {
			filtered = append(filtered, route)
		}
	}

	return filtered
}

// match check that provided value match the search or not.
func (r *labResource) match(search string, values ...string) bool {
	for _, val := range values {
		if strings.Contains(val, search) {
			return true
		}
	}
	return false
}

// Pong is like health-check route.
func (r *labResource) Pong(c echo.Context) error {
	return hecho.Write(c, RespSuccessPong)
}
