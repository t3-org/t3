package commands

import (
	"time"

	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/propagation"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/router/api"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Management of the http server",
}

var serverListenCMD = &cobra.Command{
	Use:     "listen",
	Short:   "Run the http server",
	Example: "listen",
	RunE:    withApp(serverCmdF),
}

func init() {
	serverCmd.AddCommand(serverListenCMD)

	rootCmd.AddCommand(serverCmd)
}

func tuneEcho(cfg *config.Config, sp base.ServiceProvider, a app.App) *echo.Echo {
	metricsCfg := hecho.MetricsConfig{
		MeterProvider: sp.OpenTelemetry().MeterProvider(),
	}

	tracingCfg := hecho.TracingConfig{
		Propagator:     propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		TracerProvider: sp.OpenTelemetry().TracerProvider(),
		ServerName:     cfg.ServiceName(),
	}

	e := echo.New()

	// set limit options
	e.Server.ReadTimeout = time.Millisecond * time.Duration(cfg.RequestReadTimeoutMs)
	e.Server.ReadHeaderTimeout = time.Millisecond * time.Duration(cfg.RequestReadHeaderTimeoutMs)
	e.Server.WriteTimeout = time.Millisecond * time.Duration(cfg.ResponseWriteTimeoutMs)
	e.Server.IdleTimeout = time.Millisecond * time.Duration(cfg.ConnectionIdleTimeoutMs)
	e.Server.MaxHeaderBytes = cfg.MaxHeaderSizeKb << 10 // kb to bytes

	e.HideBanner = true
	e.Logger = hecho.HexaToEchoLogger(sp.Logger(), cfg.EchoLogLevel)
	e.Debug = cfg.Debug
	e.HTTPErrorHandler = hecho.HTTPErrorHandler(sp.Logger(), sp.Translator(), e.Debug)

	//e.Use(hecho.LimitBodySize(cfg.MaxBodySizeKb << 10))

	// CORS HEADERS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           cfg.CorsMaxAgeSeconds,
	}))

	// Log requests
	e.Use(middleware.Logger())

	e.Use(hecho.Metrics(metricsCfg)) // Metrics
	e.Use(hecho.Tracing(tracingCfg)) // Distributed tracing

	// Recover recovers each panic and returns its to the echo error handler
	e.Use(hecho.Recover())

	// RequestID set requestID on each request that has blank request id.
	e.Use(hecho.RequestID())

	// CorrelationID set X-Correlation-ID value.
	e.Use(hecho.CorrelationID())

	e.Use(hecho.ExtractAuthToken(
		hecho.HeaderAuthTokenExtractor(hecho.TokenHeaderAuthorization),
		hecho.CookieTokenExtractor(cfg.AuthTokenCookie.Name),
	))

	//e.Use(middleware.CSRFWithConfig(cfg.CSRFConfig()))

	// Set user in each request context.
	//e.Use(hecho.CurrentUserBySub(a.UserFinder()))

	// HexaContext set hexa context on each request.
	e.Use(hecho.HexaContext(sp.Logger(), sp.Translator()))

	// SetContextLogger set the echo logger on each echo's context.
	e.Use(hecho.SetContextLogger(cfg.EchoLogLevel))

	// Add more data to each trace span:
	e.Use(hecho.TracingDataFromUserContext())

	return e
}

func bootEcho(cfg *config.Config, sp base.ServiceProvider, app app.App) *hecho.EchoService {
	e := tuneEcho(cfg, sp, app)

	// Register Routes
	// Register Routes
	(&api.API{
		Echo:  e,
		API:   e.Group("api/v1"),
		App:   app,
		SP:    sp,
		Guest: hecho.GuestMiddleware(),
		Auth:  hecho.AuthMiddleware(),
		Debug: hecho.DebugMiddleware(e),
	}).RegisterRoutes()

	echoService := &hecho.EchoService{Echo: e, Address: cfg.ListeningAddress()}
	registry.Register(registry.HttpServerService, echoService)

	return echoService
}

func serverCmdF(o *cmdOpts, cmd *cobra.Command, args []string) error {
	cfg := o.Cfg
	e := bootEcho(cfg, o.SP, o.App)

	if err := runProbeServer(o.SP); err != nil {
		return tracer.Trace(err)
	}

	// Start server
	app.Banner("Shield")
	return tracer.Trace(e.Run())
}
