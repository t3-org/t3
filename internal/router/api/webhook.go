package api

import (
	"github.com/kamva/hexa"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"t3.org/t3/internal/app"
	apperr "t3.org/t3/internal/error"
	"t3.org/t3/internal/input"
)

func (api *API) registerWebhookRoutes(webhooks *echo.Group, res *webhookResource) {
	webhooks.POST("/:webhook_type", res.Handle).Name = "webhooks::call"
}

type webhookResource struct {
	Resource
	app app.App
}

func (r *webhookResource) Handle(c echo.Context) error {
	// Check if the
	switch c.Param("webhook_type") {
	case "grafana":
		return r.handleGrafanaWebhook(c)
	default:
		return apperr.ErrInvalidWebhookType.SetData(hexa.Map{"reason": "invalid webhook type"})
	}
}

func (r *webhookResource) handleGrafanaWebhook(c echo.Context) error {
	var webhook input.GrafanaWebhookPayload
	if err := c.Bind(&webhook); err != nil {
		return tracer.Trace(err)
	}

	in := input.BatchUpsertTickets{Tickets: webhook.ToPatchInputs()}
	in.RemoveInternalLabels() // We do not allow API to touch internal labels:

	dto, err := r.app.UpsertTickets(c.Request().Context(), &in)
	if err != nil {
		return tracer.Trace(err)
	}

	return hecho.Write(c, RespSuccessHandleWebhook.SetData(dto))
}
