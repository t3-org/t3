package arranger

import (
	"context"
	"errors"
	"net/http"

	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"go.temporal.io/sdk/workflow"
)

// UnknownHexaErr is an error to use when we get unknown error from activity or workflow
// and need to wrap it before reporting.
var UnknownHexaErr = hexa.NewError(http.StatusInternalServerError, "lib.unknown_hexa_error")

func HandleErrWithCtx(ctx context.Context, err error) error {
	return HandleErr(err, hexa.Logger(ctx), hexatranslator.CtxTranslator(ctx))
}

func HandleErr(err error, l hlog.Logger, t hexa.Translator) error {
	ReportErr(l, err)
	return HexaToApplicationErr(err, t)
}

// ReportErr reports our error:
// supported types:
// - wrapped hexa error.
// - ApplicationError which is result of converting hexa error to ApplicationError.
// - otherwise warp error un UnknownHexaErr err.
//
// if error is hexa error, so report hexa error.
// if error is ApplicationError and if we can convert it to hexa error, so we
// convert it before report.
// otherwise wrap error in UnknownHexaErr error before report.
// TODO: we can also extend reporter to report ApplicationErrors which
// are not hexa error(e.g., CancelActivity,...) properly and in a good
// format.
func ReportErr(l hlog.Logger, e error) {
	var continueAsNew *workflow.ContinueAsNewError
	// We don't need to report nil or continueAsNew errors.
	if e == nil || errors.As(e, &continueAsNew) {
		return
	}

	hexaErr := hexa.AsHexaErr(e)
	if hexaErr == nil {
		// Maybe this is another error returned from activity and is
		// wrapped in application error.
		// (w.g., when workflow get error from activity)
		if err, ok := HexaErrFromApplicationErrWithOk(e); ok {
			hexaErr = err.(hexa.Error)
		}
	}

	if hexaErr == nil {
		hexaErr = UnknownHexaErr.SetError(tracer.Trace(e))
	}

	hexaErr.ReportIfNeeded(l, nil)
}
