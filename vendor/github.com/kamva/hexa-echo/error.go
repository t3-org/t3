package hecho

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

// HTTPErrorHandler is the echo error handler.
// This function needs to the HexaContext middleware.
func HTTPErrorHandler(l hlog.Logger, t hexa.Translator, debug bool) echo.HTTPErrorHandler {
	return func(rErr error, c echo.Context) {
		hexaErr := hexa.AsHexaErr(rErr)
		if hexaErr == nil {
			httpErr := &echo.HTTPError{}
			if errors.As(rErr, &httpErr) {
				hexaErr = errEchoHTTPError.SetHTTPStatus(httpErr.Code)
				if httpErr.Code == http.StatusNotFound {
					hexaErr = errHTTPNotFoundError
					// NOTE: Do not set the "Internal" field of the http.StatusNotFound error.
					// otherwise for next 404 requests Echo checks if its internal error field
					// is not empty, it pass the internal field to this function as error instead
					// of real 404 error !

					httpErr.Message = fmt.Sprintf("route %s %s not found", c.Request().Method, c.Request().URL)
				}

				hexaErr = hexaErr.SetError(tracer.MoveStackIfNeeded(rErr, httpErr))
			} else {
				hexaErr = errUnknownError.SetError(rErr)
			}
		}

		// Maybe error occur before set hexa context in middleware
		ctx := c.Request().Context()
		logger := hexa.CtxLogger(ctx)
		translator := hexa.CtxTranslator(ctx)

		if logger == nil {
			logger = l
		}
		if translator == nil {
			translator = t
		}

		handleError(hexaErr, c, logger, translator, debug)
	}

}

func handleError(hexaErr hexa.Error, c echo.Context, l hlog.Logger, t hexa.Translator, debug bool) {
	msg, err := hexaErr.Localize(t)

	if err != nil {
		l.With(hlog.String("translation_key", hexaErr.ID()), hlog.Err(err)).Warn("translation for error id not found.")
	}

	// Report
	hexaErr.ReportIfNeeded(l, t)

	debugData := hexaErr.ReportData()
	if debugData == nil {
		debugData = make(map[string]any, 0)
	}

	debugData["err"] = hexaErr.Error()

	body := &hexa.HTTPRespBody{
		Code:    hexaErr.ID(),
		Message: msg,
		Data:    hexaErr.Data(),
	}
	if debug {
		body.Debug = debugData
	}

	w := c.Response().Writer
	w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(hexaErr.HTTPStatus())
	_, err = easyjson.MarshalToWriter(body, w)
	if err != nil {
		l.Error("occurred error on request", hlog.Err(err))
	}
}
