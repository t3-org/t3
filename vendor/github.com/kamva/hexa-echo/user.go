package hecho

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
)

type (
	// UserFinderBySub find the user by provided sub.
	UserFinderBySub func(ctx context.Context, sub string) (hexa.User, error)

	// CurrentUserConfig is the config to use in CurrentUser middleware.
	CurrentUserConfig struct {
		uf             UserFinderBySub // Can be nil if ExtendJWT is false.
		ExtendJWT      bool
		UserContextKey string
		JWTContextKey  string
	}

	CurrentUserBySubConfig struct {
		UserFinder     UserFinderBySub
		SubContextKey  string
		UserContextKey string
	}
)

var (
	// CurrentUserContextKey is the context key to set
	// the current user in the request context.
	CurrentUserContextKey = "user"
	SubContextKey         = "sub"
)

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUser(uf UserFinderBySub) echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:      true,
		uf:             uf,
		UserContextKey: CurrentUserContextKey,
		JWTContextKey:  JwtContextKey,
	})
}

// CurrentUserWithoutFetch is for when you have a gateway that find the user and include
// it in the jwt. so you will dont need to any user finder.
func CurrentUserWithoutFetch() echo.MiddlewareFunc {
	return CurrentUserWithConfig(CurrentUserConfig{
		ExtendJWT:      false,
		uf:             nil,
		UserContextKey: CurrentUserContextKey,
		JWTContextKey:  JwtContextKey,
	})
}

// CurrentUser is a middleware to set the user in the context.
// If provided jwt, so this function find user and set it as user
// otherwise set guest user.
func CurrentUserWithConfig(cfg CurrentUserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {

			var user = hexa.NewGuest()

			// Get jwt (if exists)
			if token, ok := ctx.Get(cfg.JWTContextKey).(*jwt.Token); ok {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					if cfg.ExtendJWT {
						user, err := cfg.uf(ctx.Request().Context(), claims["sub"].(string))
						if err != nil {
							err = tracer.Trace(err)
							return err
						}
						gutil.ExtendMap(claims, user.MetaData(), true)
					}

					user, err = hexa.NewUserFromMeta(hexa.Map(claims))
					if err != nil {
						err = tracer.Trace(err)
						return
					}

				} else {
					return errors.New("JWT claims are not valid")
				}

			}

			// Set user in context with the given key
			ctx.Set(cfg.UserContextKey, user)

			// Also set for user to ua in hexa context
			ctx.Set(ContextKeyHexaUser, user)

			return next(ctx)
		}
	}
}

func CurrentUserBySub(uf UserFinderBySub) echo.MiddlewareFunc {
	return CurrentUserBySubWithConfig(CurrentUserBySubConfig{
		UserFinder:     uf,
		SubContextKey:  SubContextKey,
		UserContextKey: CurrentUserContextKey,
	})
}

func CurrentUserBySubWithConfig(cfg CurrentUserBySubConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var user = hexa.NewGuest()
			var err error
			sub, ok := ctx.Get(cfg.SubContextKey).(string)

			if ok {
				user, err = cfg.UserFinder(ctx.Request().Context(), sub)
				if err != nil {
					err = tracer.Trace(err)
					return err
				}
			}

			// Set user in context with the given key
			ctx.Set(cfg.UserContextKey, user)

			// Also set for user to use in hexa context
			ctx.Set(ContextKeyHexaUser, user)

			return next(ctx)
		}
	}
}
