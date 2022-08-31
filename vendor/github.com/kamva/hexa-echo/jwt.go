package hecho

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Cryptography algorithms to sign our jwt token
const (
	AlgorithmHS256 = "HS256"
	AlgorithmHS384 = "HS384"
	AlgorithmHS512 = "HS512"

	AlgorithmRS256 = "RS256"
	AlgorithmRS384 = "RS384"
	AlgorithmRS512 = "RS512"

	AlgorithmES256 = "ES256"
	AlgorithmES384 = "ES384"
	AlgorithmES512 = "ES512"
)

// TokenAuthorizer authorize the token.
type TokenAuthorizer func(claims jwt.MapClaims) error

type SubGenerator func(user hexa.User) (string, error)

// GenerateTokenConfig use as config to generate new token.
type GenerateTokenConfig struct {
	SingingMethod    jwt.SigningMethod
	Key              any // for rsa this is the private key
	SubGenerator     SubGenerator
	Claims           jwt.MapClaims
	ExpireTokenAfter time.Duration
}

// AuthorizeRefreshTokenConfig use as config to refresh access token.
type AuthorizeRefreshTokenConfig struct {
	SingingMethod jwt.SigningMethod
	Key           any // for rsa this is the public key
	Token         string
	// Use Authorizer to verify that can get new token.
	Authorizer TokenAuthorizer
	UserFinder UserFinderBySub
}

//--------------------------------
// JWT Middleware
//--------------------------------

const JwtContextKey = "jwt"

// SkipIfNotProvidedHeader skip jwt middleware if jwt authorization header
// is not provided.
func SkipIfNotProvidedHeader(header string) middleware.Skipper {
	return func(c echo.Context) bool {
		return c.Request().Header.Get(header) == ""
	}
}

// JwtErrorHandler check errors type and return relative hexa error.
func JwtErrorHandler(err error) error {
	// missing or malformed jwt token
	if err == middleware.ErrJWTMissing {
		return errJwtMissing.SetError(tracer.Trace(err))
	}

	// otherwise authorization error
	return errInvalidOrExpiredJwt.SetError(tracer.Trace(err))
}

func AuthorizeAudience(aud string) TokenAuthorizer {
	return func(claims jwt.MapClaims) error {
		claimAud, ok := claims["aud"]
		if !ok {
			return errInvalidAudience
		}

		if claimAud.(string) == "" {
			return errInvalidAudience
		}

		audList := strings.Split(claimAud.(string), ",")
		if !gutil.Contains(audList, aud) {
			return errInvalidAudience
		}

		return nil
	}
}

//--------------------------------
// JWT claim authorizer
//--------------------------------

type JwtClaimAuthorizerConfig struct {
	Skipper       middleware.Skipper
	JWTContextKey string
	Authorizer    TokenAuthorizer
}

func JwtClaimAuthorizer(cfg JwtClaimAuthorizerConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {

			if cfg.Skipper == nil {
				return errors.New("skipper can not be nil")
			}

			if cfg.Skipper(ctx) {
				return next(ctx)
			}

			token, ok := ctx.Get(cfg.JWTContextKey).(*jwt.Token)
			// Get jwt (if exists)
			if !ok {
				return errors.New("can not find jwt token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return errors.New("invalid jwt claims")
			}

			if err := cfg.Authorizer(claims); err != nil {
				return tracer.Trace(err)
			}

			return next(ctx)
		}
	}
}

//--------------------------------
// JWT Generator
//--------------------------------

// IDAsSubjectGenerator return user's id as jwt subject (sub).
func IDAsSubjectGenerator(user hexa.User) (string, error) {
	return user.Identifier(), nil
}

// GenerateToken generate new token for the user.
func GenerateToken(u hexa.User, cfg GenerateTokenConfig) (token string, err error) {
	if err = tracer.Trace(validateGenerateTokenCfg(cfg)); err != nil {
		return
	}

	sub, err := cfg.SubGenerator(u)
	if err != nil {
		err = tracer.Trace(err)
		return
	}
	if cfg.Claims == nil {
		cfg.Claims = make(map[string]any)
	}
	gutil.ExtendMap(cfg.Claims, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(cfg.ExpireTokenAfter).Unix(),
	}, true)

	jToken := jwt.New(cfg.SingingMethod)
	// Set claims
	jToken.Claims = cfg.Claims
	// Generate encoded token and send it as response.
	t, err := jToken.SignedString(cfg.Key)
	return t, tracer.Trace(err)
}

// AuthorizeRefreshToken authorize provided jwt token. it first fetch the user from token, and then
// provide user and claim to the authorizer.
func AuthorizeRefreshToken(cfg AuthorizeRefreshTokenConfig) (user hexa.User, err error) {
	if err = tracer.Trace(validateRefreshTokenCfg(cfg)); err != nil {
		return
	}

	// Parse token:
	jToken, err := jwt.Parse(cfg.Token, func(token *jwt.Token) (any, error) {
		return cfg.Key, nil
	})

	if err != nil {
		err = errInvalidRefreshToken.SetError(tracer.Trace(err))
		return
	}

	// Authorize user to verify user can get new access token.
	claims := jToken.Claims.(jwt.MapClaims)
	user, err = cfg.UserFinder(context.Background(), claims["sub"].(string))
	if err != nil {
		return nil, tracer.Trace(err)
	}

	return user, cfg.Authorizer(claims)
}

func validateGenerateTokenCfg(cfg GenerateTokenConfig) error {
	if gutil.IsNil(cfg.SingingMethod) || gutil.IsNil(cfg.Key) {
		return tracer.Trace(errors.New("invalid config values to generate token pairs"))
	}

	if cfg.SubGenerator == nil {
		return tracer.Trace(errors.New("invalid subject uuidGenerator for jwt"))
	}

	return nil
}

func validateRefreshTokenCfg(cfg AuthorizeRefreshTokenConfig) error {
	if cfg.Authorizer == nil {
		return tracer.Trace(errors.New("authorizer can not be nil"))
	}

	if cfg.UserFinder == nil {
		return tracer.Trace(errors.New("user finder can not be nil"))
	}

	if cfg.Token == "" {
		return errRefreshTokenCanNotBeEmpty
	}

	return nil
}
