package hecho

import (
	"errors"
	"net/http"

	"github.com/kamva/hexa"
)

var (
	// Error code description:
	// hec = hexa echo (package or project name)
	// 1 = errors about user section (identify some part in application)
	// E = Error (type of code : error|response)
	// 00 = error number zero (id of code in that part and type)

	//--------------------------------
	// User and authentication Errors
	//--------------------------------
	errUserNotFound = hexa.NewError(http.StatusInternalServerError, "lib.user.not_found_error")

	errContextUserNotImplementedHexaUser = hexa.NewError(
		http.StatusInternalServerError,
		"lib.user.interface_not_implemented_error",
	).SetError(errors.New("user in the hexa context does not implemented User interface"))

	errJwtMissing = hexa.NewError(
		http.StatusBadRequest,
		"lib.user.missing_jwt_token_error",
	).SetError(errors.New("missing or malformed jwt"))

	errInvalidOrExpiredJwt = hexa.NewError(
		http.StatusUnauthorized,
		"lib.user.invalid_expired_jwt_error",
	).SetError(errors.New("invalid or expired jwt"))

	errInvalidAudience = hexa.NewError(
		http.StatusUnauthorized,
		"lib.user.invalid_jwt_audience_error",
	).SetError(errors.New("audience value in the jwt token is not for this app"))

	errUserMustBeGuest = hexa.NewError(http.StatusUnauthorized, "lib.user.must_be_guest_error")

	errUserNeedToAuthenticate = hexa.NewError(http.StatusUnauthorized, "lib.user.must_authenticate_error")

	errRefreshTokenCanNotBeEmpty = hexa.NewError(http.StatusBadRequest, "lib.user.refresh_token_is_empty_error")

	errInvalidRefreshToken = hexa.NewError(http.StatusBadRequest, "lib.user.invalid_refresh_token_error")

	//--------------------------------
	// Request errors
	//--------------------------------
	errRequestIdNotFound = hexa.NewError(http.StatusInternalServerError, "lib.request.request_id_not_found_error")

	errCorrelationIDNotFound = hexa.NewError(http.StatusInternalServerError, "lib.request.correlation_id_not_found_error")

	//--------------------------------
	// DEBUG
	//--------------------------------
	errRouteAvailableInDebugMode = hexa.NewError(http.StatusUnauthorized, "lib.route.available_in_debug_mode_error")

	//--------------------------------
	// Other errors
	//--------------------------------
	errHTTPNotFoundError = hexa.NewError(http.StatusNotFound, "lib.route.not_found_error")

	// Set this error status manually on return relative to echo error code.
	errEchoHTTPError = hexa.NewError(http.StatusNotFound, "lib.http_server.occurred_http_error")

	errUnknownError = hexa.NewError(http.StatusInternalServerError, "lib.http_server.unknown_error")
)
