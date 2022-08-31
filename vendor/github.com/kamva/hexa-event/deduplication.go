package hevent

import "go.opentelemetry.io/otel/attribute"

const (
	// HexaEventID is id of the event we use as key in hexa Context.
	HexaEventID = "HEXA_EVENT_ID"
	// HexaEventHandlerActionName is action's name that event handler wants to do.
	HexaEventHandlerActionName = "HEXA_EVENT_HANDLER_ACTION_NAME"
	// HexaRootEventID is id of the root event for retry events.
	HexaRootEventID = "HEXA_ROOT_EVENT_ID"
	// HexaRootEventHandlerActionName is action's name of the root event for retry events.
	HexaRootEventHandlerActionName = "HEXA_ROOT_EVENT_HANDLER_ACTION_NAME"
)

// open telemetry attribute keys:
var (
	MessagingActionName     = attribute.Key("messaging.action.name")
	MessagingRootActionName = attribute.Key("messaging.action.root.name")
)
