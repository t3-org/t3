package hecho

import (
	"errors"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
)

// Write writes reply as response.
// You MUST have hexa context in your echo context to use
// this function to use its logger and translator.
func Write(c echo.Context, reply hexa.Reply) error {
	ctx := c.Request().Context()
	if ctx == nil {
		return tracer.Trace(errors.New("invalid hexa context, we can not write reply as response"))
	}
	l := hexa.Logger(ctx)
	t := hexatranslator.CtxTranslator(ctx)

	return WriteWithOpts(c, l, t, reply)
}

// WriteWithOpts writes the reply as response.
func WriteWithOpts(c echo.Context, l hlog.Logger, t hexa.Translator, reply hexa.Reply) error {
	msg, err := t.Translate(reply.ID())
	if err != nil {
		l.With(hlog.String("translation_key", reply.ID())).Warn("translation for reply id not found.")
	}

	body := &hexa.HTTPRespBody{
		Code:    reply.ID(),
		Message: msg,
		Data:    reply.Data(),
	}

	w := c.Response().Writer
	w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(reply.HTTPStatus())
	_, err = easyjson.MarshalToWriter(body, w)
	if err != nil {
		l.Error("occurred error on request", hlog.Err(err))
	}
	return tracer.Trace(err)
}
