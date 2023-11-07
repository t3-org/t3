package apperr

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Base Errors
//--------------------------------

var (
	ErrInvalidIDValue = hexa.NewError(http.StatusBadRequest, "itrack.invalid_id_value")
)

//--------------------------------
// Ticket errors
//--------------------------------

var (
	ErrTicketNotFound = hexa.NewError(http.StatusNotFound, "itrack.ticket.not_found_error")
)

// --------------------------------
// Ticket KeyValue error
// --------------------------------

var (
	ErrTicketKVNotFound = hexa.NewError(http.StatusNotFound, "itrack.ticket_kv.not_found_error")
)

//--------------------------------
// System errors
//--------------------------------

var (
	ErrSystemPropertyNotFound = hexa.NewError(http.StatusNotFound, "itrack.system.property_not_found")
)

//--------------------------------
// Gateway error
//--------------------------------

var (
	ErrTooManyRequests = hexa.NewError(http.StatusTooManyRequests, "itrack.gateway.too_many_requests")
)

// --------------------------------
// Webhook errors
// --------------------------------
var (
	ErrInvalidWebhookType = hexa.NewError(http.StatusBadRequest, "itrack.webhook.err")
)
