package hecho

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CSRFSkipperByAuthTokenLocation skips if request doesn't need to csrf check.
// We do csrf check when user's token is in the cookie or session.
// and the request method is post too.
func CSRFSkipperByAuthTokenLocation(ctx echo.Context) bool {
	l, ok := ctx.Get(AuthTokenLocationContextKey).(TokenLocation)
	return !(ok && (l == TokenLocationCookie || l == TokenLocationSession))
}

var _ middleware.Skipper = CSRFSkipperByAuthTokenLocation
