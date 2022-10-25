package apperr

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Base Errors
//--------------------------------

var (
	ErrInvalidIDValue = hexa.NewError(http.StatusBadRequest, "space.invalid_id_value")
)

//--------------------------------
// Planet errors
//--------------------------------

var (
	ErrPlanetNotFound = hexa.NewError(http.StatusNotFound, "space.planet.not_found_error")
)

//--------------------------------
// System errors
//--------------------------------

var (
	ErrSystemPropertyNotFound = hexa.NewError(http.StatusNotFound, "space.system.property_not_found")
)

//--------------------------------
// Gateway error
//--------------------------------

var (
	ErrTooManyRequests = hexa.NewError(http.StatusTooManyRequests, "space.gateway.too_many_requests")
)
