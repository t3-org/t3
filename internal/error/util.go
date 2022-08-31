package apperr

import "github.com/kamva/hexa"

func IsInternalErr(err error) bool {
	if err == nil {
		return false
	}
	if herr := hexa.AsHexaErr(err); herr != nil {
		return herr.HTTPStatus() >= 500
	}
	return true // unknown error, so its internal.
}
