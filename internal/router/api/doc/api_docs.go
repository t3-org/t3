// Package doc.
//
// # T3 API docs
//
// Terms Of Service:
//
//	Schemes: http, https
//	Host: t3.app
//	BasePath:
//	Version: 0.1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- bearerAuth:
//
//	SecurityDefinitions:
//	bearerAuth:
//	     type: apiKey
//	     name: Authorization
//	     in: header
//
// swagger:meta
package doc

import (
	"github.com/kamva/hexa/pagination"
	"github.com/labstack/echo/v4"
	"t3.org/t3/internal/dto"
	"t3.org/t3/internal/input"
)

type replyCode struct {
	// response code
	Code string `json:"code"`
}

// route:begin: lab::ping
// swagger:route GET /api/v1/lab/ping lab labPingParams
// Ping API server.
// responses:
//   200: labPingSuccessResponse

// swagger:parameters labPingParams
type labPingParamsWrapper struct {
	// in:body
	Body struct {
	}
}

// success response
// swagger:response labPingSuccessResponse
type labPingResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
	}
}

// route:end: lab::ping

// route:begin: lab::routes
// swagger:route GET /api/v1/lab/routes lab labRoutesParams
// Get all routes(Debug mode).
// responses:
//   200: labRoutesSuccessResponse

// swagger:parameters labRoutesParams
type labRoutesParamsWrapper struct {
	// in:body
	Body struct {
	}
}

// success response
// swagger:response labRoutesSuccessResponse
type labRoutesResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data struct {
			Routes []*echo.Route `json:"routes"`
		} `json:"data"`
	}
}

// route:end: lab::routes

// route:begin: tickets::create
// swagger:route POST /api/v1/tickets tickets ticketsCreateParams
// Create a ticket.
// responses:
//   200: ticketsCreateSuccessResponse

// swagger:parameters ticketsCreateParams
type ticketsCreateParamsWrapper struct {
	// in:body
	Body struct {
		input.CreateTicket
	}
}

// success response
// swagger:response ticketsCreateSuccessResponse
type ticketsCreateResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Ticket `json:"data"`
	}
}

// route:end: tickets::create

// route:begin: tickets::delete
// swagger:route DELETE /api/v1/tickets/{id} tickets ticketsDeleteParams
// Delete a ticket.
// responses:
//   200: ticketsDeleteSuccessResponse

// swagger:parameters ticketsDeleteParams
type ticketsDeleteParamsWrapper struct {
	// in:path
	Id string `json:"id"`

	// in:body
	Body struct {
	}
}

// success response
// swagger:response ticketsDeleteSuccessResponse
type ticketsDeleteResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
	}
}

// route:end: tickets::delete

// route:begin: tickets::get
// swagger:route GET /api/v1/tickets/{id} tickets ticketsGetParams
//
// responses:
//   200: ticketsGetSuccessResponse

// swagger:parameters ticketsGetParams
type ticketsGetParamsWrapper struct {
	// in:path
	Id string `json:"id"`

	// in:body
	Body struct {
	}
}

// success response
// swagger:response ticketsGetSuccessResponse
type ticketsGetResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Ticket `json:"data"`
	}
}

// route:end: tickets::get

// route:begin: tickets::patch
// swagger:route PATCH /api/v1/tickets/{id} tickets ticketsPatchParams
// Patch a ticket.
// responses:
//   200: ticketsPatchSuccessResponse

// swagger:parameters ticketsPatchParams
type ticketsPatchParamsWrapper struct {
	// in:path
	Id string `json:"id"`

	// in:body
	Body struct {
		input.PatchTicket
	}
}

// success response
// swagger:response ticketsPatchSuccessResponse
type ticketsPatchResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Ticket `json:"data"`
	}
}

// route:end: tickets::patch

// route:begin: tickets::query
// swagger:route GET /api/v1/tickets tickets ticketsQueryParams
// Query tickets.
// responses:
//   200: ticketsQuerySuccessResponse

// swagger:parameters ticketsQueryParams
type ticketsQueryParamsWrapper struct {
	// in:body
	Body struct {
		// Query should be in k8s label-selector format. read its
		// docs on T3 dashbaord in the tickets search page.
		Query string `json:"query"`
		Pagination
	}
}

// success response
// swagger:response ticketsQuerySuccessResponse
type ticketsQueryResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data struct {
			pagination.Pagination
			Ietms []dto.Ticket `json:"ietms"`
		}
	}
}

// route:end: tickets::query

// route:begin: webhooks::call
// swagger:route POST /api/v1/webhooks/{webhook_type} webhooks webhooksCallParams
// Webhook endpoint.
// This Endpoint is called by sources(grafana...) as the webhook endpoint.
// responses:
//   200: webhooksCallSuccessResponse

// swagger:parameters webhooksCallParams
type webhooksCallParamsWrapper struct {
	// Its value could be : grafana
	// in:path
	Webhook_type string `json:"webhook_type"`

	// in:body
	input.GrafanaWebhookPayload

	// in:body
	Body struct {
		// Body is variable based on the webhook type. for example
		// for Grafana its the Grafana alert data.
		input.GrafanaWebhookPayload
	}
}

// success response
// swagger:response webhooksCallSuccessResponse
type webhooksCallResponseWrapper struct {
	// in: body
	Body struct {
		replyCode
		Data dto.Ticket `json:"data"`
	}
}

// route:end: webhooks::call
