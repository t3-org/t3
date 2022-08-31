package hecho

import (
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

const (

	// ContextKeyHexaRequestID uses as key in context to store request id to use in context middleware
	ContextKeyHexaRequestID = "_hexa_ctx.rid"

	// ContextKeyHexaCorrelationID uses as key in context to store correlation id to use in context middleware
	ContextKeyHexaCorrelationID = "_hexa_ctx.cid"

	// ContextKeyHexaUser is the identifier to set the hexa user as a field in the context of a request.
	ContextKeyHexaUser = "_hexa_ctx.user"
)

// getHexaUser returns hexa user instance from the current user.
func getHexaUser(ctx echo.Context) (hexa.User, hexa.Error) {
	// Get user if exists:
	u := ctx.Get(ContextKeyHexaUser)

	if u == nil {
		return nil, errUserNotFound
	}

	if u, ok := u.(hexa.User); ok {
		return u, nil
	} else {
		return nil, errContextUserNotImplementedHexaUser
	}
}

// getCorrelationID returns the request correlation id.
func getCorrelationID(ctx echo.Context) (string, hexa.Error) {
	// Get Request ID if exists:
	cid := ctx.Get(ContextKeyHexaCorrelationID).(string)

	if cid == "" {
		return "", errCorrelationIDNotFound
	}

	return cid, nil
}

// HexaContext set hexa context on each request.
func HexaContext(logger hlog.Logger, translator hexa.Translator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r := ctx.Request()
			user, err := getHexaUser(ctx)
			if err != nil {
				return tracer.Trace(err)
			}

			cid, err := getCorrelationID(ctx)
			if err != nil {
				return tracer.Trace(err)
			}

			// Set context
			hexaCtx := hexa.NewContext(r.Context(), hexa.ContextParams{
				Request:        r,
				CorrelationId:  cid,
				Locale:         r.Header.Get("Accept-Language"),
				User:           user,
				BaseLogger:     logger,
				BaseTranslator: translator,
			})

			ctx.SetRequest(r.Clone(hexaCtx))
			return next(ctx)
		}
	}
}
