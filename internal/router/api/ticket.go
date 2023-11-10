package api

import (
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"space.org/space/internal/app"
	"space.org/space/internal/input"
)

func (api *API) registerTicketRoutes(tickets *echo.Group, res *ticketResource) {
	hecho.ResourceAPI(tickets, res, "tickets")
}

type ticketResource struct {
	Resource
	app app.App
}

func (r *ticketResource) Get(c echo.Context) error {
	t, err := r.app.GetTicket(c.Request().Context(), r.ID(c))
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessGetTicket.SetData(t))
}

func (r *ticketResource) Create(c echo.Context) error {
	var in input.CreateTicket
	if err := c.Bind(&in); err != nil {
		return tracer.Trace(err)
	}

	// We do not allow API to touch internal labels:
	input.RemoveInternalLabels(in.Labels)

	dto, err := r.app.CreateTicket(c.Request().Context(), &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessCreateTicket.SetData(dto))
}

func (r *ticketResource) Patch(c echo.Context) error {
	var in input.PatchTicket
	if err := c.Bind(&in); err != nil {
		return tracer.Trace(err)
	}

	// We do not allow API to touch internal labels:
	input.RemoveInternalLabels(in.Labels)

	dto, err := r.app.PatchTicket(c.Request().Context(), r.ID(c), &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessUpdateTicket.SetData(dto))
}

func (r *ticketResource) Delete(c echo.Context) error {
	err := r.app.DeleteTicket(c.Request().Context(), r.ID(c))
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessDeleteTicket)
}

func (r *ticketResource) Query(c echo.Context) error {
	_, page, perPage := r.QueryAndPaginationParams(c)
	res, err := r.app.QueryTickets(c.Request().Context(), c.QueryParam("query"), page, perPage)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessQueryTicket.SetData(res))
}

var _ hecho.CreateResource = &ticketResource{}
var _ hecho.GetResource = &ticketResource{}
var _ hecho.PatchResource = &ticketResource{}
var _ hecho.DeleteResource = &ticketResource{}
var _ hecho.QueryResource = &ticketResource{}
