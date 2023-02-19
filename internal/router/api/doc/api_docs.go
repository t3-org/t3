// Package doc.
//
// Space API docs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: space.app
//     BasePath:
//     Version: 0.1.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearerAuth:
//
//     SecurityDefinitions:
//     bearerAuth:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package doc

import (
	"github.com/kamva/hexa/pagination"
	"github.com/labstack/echo/v4"
	"space.org/space/internal/dto"
	"space.org/space/internal/input"
)

type replyCode struct {
	// response code
	Code string `json:"code"`
}

// route:begin: lab::ping
// swagger:route GET /api/v1/lab/ping lab labPingParams
//
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
// Returns routes.
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

// route:begin: planets::create
// swagger:route POST /api/v1/planets planets planetsCreateParams
// Create a planet.
// responses:
//   200: planetsCreateSuccessResponse

// swagger:parameters planetsCreateParams
type planetsCreateParamsWrapper struct {
	// in:body
	Body struct {
		input.CreatePlanet
	}
}

// success response
// swagger:response planetsCreateSuccessResponse
type planetsCreateResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Planet `json:"data"`
	}
}

// route:end: planets::create

// route:begin: planets::delete
// swagger:route DELETE /api/v1/planets/{id} planets planetsDeleteParams
// Delete a planet.
// responses:
//   200: planetsDeleteSuccessResponse

// swagger:parameters planetsDeleteParams
type planetsDeleteParamsWrapper struct {
	// in:path
	Id string `json:"id"`

	// in:body
	Body struct {
	}
}

// success response
// swagger:response planetsDeleteSuccessResponse
type planetsDeleteResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
	}
}

// route:end: planets::delete

// route:begin: planets::getByCode
// swagger:route GET /api/v1/planets/code/{code} planets planetsGetByCodeParams
// Get a planet by code.
// responses:
//   200: planetsGetByCodeSuccessResponse

// swagger:parameters planetsGetByCodeParams
type planetsGetByCodeParamsWrapper struct {
	// in:path
	Code string `json:"code"`

	// in:body
	Body struct {
	}
}

// success response
// swagger:response planetsGetByCodeSuccessResponse
type planetsGetByCodeResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Planet `json:"data"`
	}
}

// route:end: planets::getByCode

// route:begin: planets::put
// swagger:route PUT /api/v1/planets/{id} planets planetsPutParams
// Update a planet.
// responses:
//   200: planetsPutSuccessResponse

// swagger:parameters planetsPutParams
type planetsPutParamsWrapper struct {
	// in:path
	Id string `json:"id"`

	// in:body
	Body struct {
		input.UpdatePlanet
	}
}

// success response
// swagger:response planetsPutSuccessResponse
type planetsPutResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data dto.Planet `json:"data"`
	}
}

// route:end: planets::put

// route:begin: planets::query
// swagger:route GET /api/v1/planets planets planetsQueryParams
// Query planets.
// responses:
//   200: planetsQuerySuccessResponse

// swagger:parameters planetsQueryParams
type planetsQueryParamsWrapper struct {
	// in:body
	Body struct {
		Query string `json:"query"`
	}
}

// success response
// swagger:response planetsQuerySuccessResponse
type planetsQueryResponseWrapper struct {
	// in:body
	Body struct {
		replyCode
		Data struct {
			pagination.Pages
			Items []*dto.Planet `json:"items"`
		} `json:"data"`
	}
}

// route:end: planets::query
