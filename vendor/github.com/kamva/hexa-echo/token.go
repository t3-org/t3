package hecho

import (
	"strings"

	"github.com/gorilla/sessions"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const TokenHeaderAuthorization = "Authorization"
const TokenCookieFieldAuthToken = "hexa_auth_token"
const TokenSessionFieldToken = "token"

const AuthTokenContextKey = "auth_token"
const AuthTokenLocationContextKey = "auth_token_location"

type TokenLocation int

const (
	TokenLocationUnknown = iota
	TokenLocationHeader
	TokenLocationCookie
	TokenLocationSession
)

// TokenExtractor extracts a token from somewhere and then returns the
// token and its location.
type TokenExtractor func(ctx echo.Context) (string, TokenLocation, error)

type ExtractTokenConfig struct {
	Skipper                 middleware.Skipper
	TokenContextKey         string
	TokenLocationContextKey string
	Extractors              []TokenExtractor
}

// ExtractAuthToken extracts the authentication token.
// If you want to ignore session, set nil store and empty session name.
func ExtractAuthToken(extractors ...TokenExtractor) echo.MiddlewareFunc {
	return ExtractTokenWithConfig(ExtractTokenConfig{
		TokenContextKey:         AuthTokenContextKey,
		TokenLocationContextKey: AuthTokenLocationContextKey,
		Extractors:              extractors,
	})
}

// ExtractTokenWithConfig extracts the authentication token from the cookie or Authorization header.
func ExtractTokenWithConfig(cfg ExtractTokenConfig) echo.MiddlewareFunc {
	if cfg.Skipper == nil {
		cfg.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			for _, e := range cfg.Extractors {
				token, location, err := e(ctx)
				if err != nil {
					return tracer.Trace(err)
				}
				if token != "" {
					ctx.Set(AuthTokenContextKey, token)
					ctx.Set(AuthTokenLocationContextKey, location)
					return next(ctx)
				}
			}

			return next(ctx)
		}
	}
}

func HeaderAuthTokenExtractor(headerFieldName string) TokenExtractor {
	return func(ctx echo.Context) (string, TokenLocation, error) {
		headerVal := ctx.Request().Header.Get(headerFieldName)
		if headerVal == "" {
			return "", 0, nil
		}

		var token string
		if len(headerVal) > 6 && strings.ToUpper(headerVal[0:6]) == "BEARER" {
			token = headerVal[7:]
		}

		if len(headerVal) > 5 && strings.ToUpper(headerVal[0:5]) == "TOKEN" {
			token = headerVal[6:]
		}

		return token, TokenLocationHeader, nil
	}
}

func CookieTokenExtractor(cookieFieldName string) TokenExtractor {
	return func(ctx echo.Context) (string, TokenLocation, error) {
		if cookie, err := ctx.Cookie(cookieFieldName); err == nil {
			return cookie.Value, TokenLocationCookie, nil
		}
		return "", 0, nil
	}
}

func SessionTokenExtractor(store sessions.Store, sessionName string, tokenField string) TokenExtractor {
	return func(ctx echo.Context) (string, TokenLocation, error) {
		if sess, err := store.Get(ctx.Request(), sessionName); sess != nil && err == nil {
			token, _ := sess.Values[tokenField].(string)
			return token, TokenLocationSession, nil
		}
		return "", 0, nil
	}
}
