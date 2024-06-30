package apperr

import (
	"net/http"

	"github.com/kamva/hexa"
)

//--------------------------------
// Base Errors
//--------------------------------

var (
	ErrInvalidIDValue = hexa.NewError(http.StatusBadRequest, "t3.invalid_id_value")
	ErrInvalidQuery   = hexa.NewError(http.StatusBadRequest, "t3.invalid_query")
)

//--------------------------------
// Ticket errors
//--------------------------------

var (
	ErrTicketNotFound              = hexa.NewError(http.StatusNotFound, "t3.ticket.not_found_error")
	ErrTicketRequiredFieldsMissing = hexa.NewError(http.StatusForbidden, "t3.ticket.required_fields_missing")
)

// --------------------------------
// Ticket KeyValue error
// --------------------------------

var (
	ErrTicketKVNotFound = hexa.NewError(http.StatusNotFound, "t3.ticket_kv.not_found_error")
)

//--------------------------------
// System errors
//--------------------------------

var (
	ErrSystemPropertyNotFound = hexa.NewError(http.StatusNotFound, "t3.system.property_not_found")
)

//--------------------------------
// Gateway error
//--------------------------------

var (
	ErrTooManyRequests = hexa.NewError(http.StatusTooManyRequests, "t3.gateway.too_many_requests")
)

// --------------------------------
// Webhook errors
// --------------------------------

var (
	ErrInvalidWebhookType = hexa.NewError(http.StatusBadRequest, "t3.webhook.err")
)
