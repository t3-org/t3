package arranger

import (
	"encoding/json"
	"errors"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"go.temporal.io/sdk/temporal"
)

// Our special error type to check in all microservices in our converts from temporal
// Application error to hexa error.
const HexaErrType = "_hexa_err"

func HexaErrFromApplicationErr(e error) error {
	e, _ = HexaErrFromApplicationErrWithOk(e)
	return e
}

// HexaErrFromApplicationErrWithOk converts provided error to hexa error if the provided error
// is(itself or one of its inner errors) an temporal.ApplicationError and also its error
// type prop is HexaErrType.
// You can use this convert function in both workflow
// and activity errors.
func HexaErrFromApplicationErrWithOk(e error) (error, bool) {
	if e == nil {
		return nil, false
	}

	var applicationErr *temporal.ApplicationError
	if !errors.As(e, &applicationErr) {
		return e, false
	}

	if applicationErr.Type() != HexaErrType {
		return e, false
	}

	var details ErrorDetails
	if err := applicationErr.Details(&details); err != nil {
		errMsg := "invalid error details, can not decode hexa error details from temporal error"
		hlog.Error(errMsg, hlog.Err(tracer.Trace(err)))
		return e, false
	}

	return errDetailsToHexaErr(&details), true
}

// HexaToApplicationErr converts provided error to *temporal.ApplicationError
// if its is hexa error. otherwise return it untouched to handle by temporal
// itself. You can use this convert function in both activity and workflow
// error.
func HexaToApplicationErr(e error, t hexa.Translator) error {
	if e == nil {
		return nil
	}

	hErr := hexa.AsHexaErr(e)
	if hErr == nil {
		// We do not convert unknown errors because maybe
		// its our activity returned a temporal error. e.g.,
		// CancelError,... or even ApplicationError which
		// is result of converting hexa error to ApplicationError
		// in the activity. so we should not convert it.
		return e
	}

	details := hexaErrToErrDetails(hErr, t)
	return temporal.NewApplicationErrorWithCause(hErr.Error(), HexaErrType, hErr.InternalError(), details)
}

func hexaErrToErrDetails(hErr hexa.Error, t hexa.Translator) *ErrorDetails {
	localMsg, _ := hErr.Localize(t)

	data, err := json.Marshal(hErr.Data())
	if err != nil {
		errMsg := "can not marshal error data"
		hlog.Error(errMsg, hlog.Err(tracer.Trace(err)))
	}

	return &ErrorDetails{
		Status:           int32(hErr.HTTPStatus()),
		Id:               hErr.ID(),
		LocalizedMessage: localMsg,
		Data:             data,
	}
}

func errDetailsToHexaErr(d *ErrorDetails) hexa.Error {
	m := make(map[string]any)
	if err := gutil.Unmarshal(d.Data, &m); err != nil {
		errMsg := "can not unmarshal error data"
		hlog.Error(errMsg, hlog.Err(tracer.Trace(err)))
	}

	hexaErr := hexa.NewLocalizedError(int(d.Status), d.Id, d.LocalizedMessage, nil)
	hexaErr = hexaErr.SetData(m)

	return hexaErr
}

