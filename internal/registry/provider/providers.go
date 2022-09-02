package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hibiken/asynq"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hcache "github.com/kamva/hexa-cache"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa-job/hsynq"
	hexarobfig "github.com/kamva/hexa-job/robfig"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
	"space.org/space/internal/model"
	"space.org/space/internal/registry"
	"space.org/space/internal/router/api"
	"space.org/space/internal/router/crons"
	"space.org/space/internal/router/jobs"
	"space.org/space/internal/store"
	cachestore "space.org/space/internal/store/cache"
	"space.org/space/internal/store/sqlstore"
	"space.org/space/pkg/hredis"
)

var providers = map[string]registry.Provider{
	registry.ServiceNameServiceProvider: ServiceProviderProvider,
	registry.ServiceNameLogger:          LoggerProvider,
	registry.ServiceNameTranslator:      TranslatorProvider,
	registry.ServiceNameProbeServer:     ProbeServerProvider,
	registry.ServiceNameHealthReporter:  HealthReporterProvider,
	registry.ServiceNameTracerProvider:  TracerProvider,
	registry.ServiceNameMeterProvider:   MeterProvider,
	registry.ServiceNameOpenTelemetry:   OpenTelemetryProvider,
	registry.ServiceNameRedis:           RedisProvider,
	registry.ServiceNameCacheProvider:   CacheProvider,
	registry.ServiceNameHttpServer:      HttpServerProvider,
	registry.ServiceNameJobs:            JobsProvider,
	registry.ServiceNameWorker:          WorkerProvider,
	registry.ServiceNameCron:            CronProvider,
	registry.ServiceNameStore:           StoreProvider,
	registry.ServiceNameApp:             AppProvider,
}

func Providers(names []string) map[string]registry.Provider {
	m := make(map[string]registry.Provider)
	for _, name := range names {
		p, ok := providers[name]
		if !ok {
			hlog.Error("can not find provider, ignoring it.", hlog.String("name", name))
			continue
		}
		m[name] = p
	}
	return m
}

func ServiceProviderProvider(r hexa.ServiceRegistry, _ *config.Config) error {
	r.Register(registry.ServiceNameServiceProvider, base.NewServiceProvider(r))
	return nil
}

func LoggerProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	l, err := logdriver.NewStackLoggerDriver(cfg.StackLoggerOptions())
	if err != nil {
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameLogger, l)
	hlog.SetGlobalLogger(l)
	return nil
}

func TranslatorProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	translator := huner.NewTranslator(cfg.I18nPath(), cfg.TranslateOptions())
	hexatranslator.SetGlobal(translator)
	r.Register(registry.ServiceNameTranslator, translator)
	return nil
}

func HealthReporterProvider(r hexa.ServiceRegistry, _ *config.Config) error {
	r.Register(registry.ServiceNameHealthReporter, hexa.NewHealthReporter())
	return nil
}

func ProbeServerProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	reporter := r.Service(registry.ServiceNameHealthReporter).(hexa.HealthReporter)
	probeServer := probe.NewServer(&http.Server{Addr: cfg.ProbeServerAddress}, http.NewServeMux())

	probe.RegisterHealthHandlers(probeServer, reporter)
	// Register other probe server handlers here.

	r.Register(registry.ServiceNameProbeServer, probeServer)
	return nil
}

func JobsProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	l := r.Service(registry.ServiceNameLogger).(hlog.Logger)
	t := r.Service(registry.ServiceNameTranslator).(hexa.Translator)

	propagator := hexa.NewContextPropagator(l, t)
	jobs := hsynq.NewJobs(asynq.NewClient(cfg.AsynqConfig.RedisOpts()), propagator, hsynq.NewJsonTransformer())

	r.Register(registry.ServiceNameJobs, jobs)
	return nil
}

func TracerProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	tcfg := cfg.OpenTelemetry.Tracing
	if tcfg.NoopTracer {
		tp := trace.NewNoopTracerProvider()
		r.Register(registry.ServiceNameTracerProvider, tp)
		return nil
	}

	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(tcfg.JaegerAddr)))
	if err != nil {
		return tracer.Trace(err)
	}

	sampler := tracesdk.AlwaysSample()
	if !cfg.Debug && !tcfg.AlwaysSample {
		// We use the ParentBased(AlwaysSample)	sampler in production.
		sampler = tracesdk.ParentBased(tracesdk.AlwaysSample())
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(sampler),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName()),
			attribute.String("environment", cfg.Environment),
			attribute.String("service_instance", cfg.InstanceName),
		)),
	)

	r.Register(registry.ServiceNameTracerProvider, tp)
	return nil
}

func MeterProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	mcfg := cfg.OpenTelemetry.Metric

	if !mcfg.Enabled {
		r.Register(registry.ServiceNameMeterProvider, metric.NewNoopMeterProvider())
		return nil
	}

	// Initialize prometheus exporter
	promCfg := prometheus.Config{DefaultHistogramBoundaries: mcfg.Prometheus.DefaultHistogramBoundaries}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(histogram.WithExplicitBoundaries(promCfg.DefaultHistogramBoundaries)),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
		controller.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName()),
			attribute.String("service_instance", cfg.InstanceName),
		)),
	)

	exporter, err := prometheus.New(promCfg, c)
	if err != nil {
		return tracer.Trace(err)
	}

	// Register probe handler
	probeserver := r.Service(registry.ServiceNameProbeServer).(probe.Server)
	probeserver.Register("prometheus_metrics", "/prometheus/metrics", exporter.ServeHTTP, "Prometheus metrics")

	// Register prometheus exporter as another service to just keep it:
	r.Register(registry.ServiceNamePrometheus, exporter)
	r.Register(registry.ServiceNameMeterProvider, exporter.MeterProvider())
	return nil
}

func OpenTelemetryProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	tp := r.Service(registry.ServiceNameTracerProvider).(trace.TracerProvider)
	mp := r.Service(registry.ServiceNameMeterProvider).(metric.MeterProvider)
	r.Register(registry.ServiceNameOpenTelemetry, htel.NewOpenTelemetry(tp, mp))
	return nil
}

func RedisProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	dial := func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", cfg.RedisAddress)
		if err != nil {
			return nil, err
		}

		if cfg.RedisPassword != "" {
			if _, err := c.Do("AUTH", cfg.RedisPassword); err != nil {
				c.Close()
				return nil, err
			}
		}

		if _, err := c.Do("SELECT", cfg.RedisDB); err != nil {
			c.Close()
			return nil, err
		}
		return c, nil
	}

	testOnBorrow := func(c redis.Conn, t time.Time) error {
		if time.Since(t) < time.Minute {
			return nil
		}
		_, err := c.Do("PING")
		return err
	}

	red := &hredis.Service{
		Pool: &redis.Pool{
			Dial:         dial,
			TestOnBorrow: testOnBorrow,
			MaxIdle:      5,
			IdleTimeout:  time.Second * 120,
		},
	}

	r.Register(registry.ServiceNameRedis, red)
	return nil
}

func CacheProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	if cfg.Cache.Enabled {
		return nil
	}

	red := r.Service(registry.ServiceNameRedis).(*hredis.Service)
	hcache.NewRedisCacheProvider(&hcache.RedisOptions{
		Prefix:      cfg.Cache.Redis.Prefix,
		Pool:        red.Pool,
		Marshaler:   hcache.GobMarshal,
		Unmarshaler: hcache.GobUnmarshal,
		DefaultTTL:  cfg.Cache.Redis.TTL(),
	})
	return nil
}

func StoreProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	sp := r.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
	var s model.Store

	s, err := sqlstore.New(sp, cfg.DB)
	if err != nil {
		hlog.Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return tracer.Trace(err)
	}

	if sp.CacheProvider() != nil { // Add the cache layer.
		s = cachestore.New(sp, s)
	}

	s = store.NewTracingLayerStore("sql", sp.OpenTelemetry().TracerProvider(), s)
	r.Register(registry.ServiceNameStore, s)

	// Set global DB store on the model package:
	model.SetDBStore(s)

	return nil
}

func AppProvider(r hexa.ServiceRegistry, _ *config.Config) error {
	sp := r.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
	s := r.Service(registry.ServiceNameStore).(model.Store)

	a, err := app.NewWithAllLayers(sp, s)
	if err != nil {
		sp.Logger().Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameApp, a)
	return nil
}

func CronProvider(r hexa.ServiceRegistry, _ *config.Config) error {
	sp := r.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
	a := r.Service(registry.ServiceNameApp).(app.App)

	u, err := hexa.NewGuest().SetMeta(hexa.UserMetaKeyName, "cron_job")
	if err != nil {
		return tracer.Trace(err)
	}

	ctxGen := func(c context.Context) context.Context {
		return hexa.NewContext(c, hexa.ContextParams{
			Locale:         "en",
			User:           u,
			CorrelationId:  gutil.UUID(),
			BaseLogger:     sp.Logger(),
			BaseTranslator: sp.Translator(),
		})
	}

	cronJobs := hexarobfig.New(ctxGen, cron.New())

	if err := crons.RegisterCronJobs(cronJobs, sp, a); err != nil {
		return tracer.Trace(err)
	}

	registry.Register(registry.ServiceNameCron, cronJobs)
	return nil
}

func tuneEcho(cfg *config.Config, sp base.ServiceProvider) *echo.Echo {
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

func HttpServerProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	sp := r.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
	a := r.Service(registry.ServiceNameApp).(app.App)

	e := tuneEcho(cfg, sp)

	// Register Routes
	(&api.API{
		Echo:  e,
		API:   e.Group("api/v1"),
		App:   a,
		SP:    sp,
		Guest: hecho.GuestMiddleware(),
		Auth:  hecho.AuthMiddleware(),
		Debug: hecho.DebugMiddleware(e),
	}).RegisterRoutes()

	echoService := &hecho.EchoService{Echo: e, Address: cfg.ListeningAddress()}
	registry.Register(registry.ServiceNameHttpServer, echoService)
	return nil
}

func WorkerProvider(r hexa.ServiceRegistry, cfg *config.Config) error {
	sp := r.Service(registry.ServiceNameServiceProvider).(base.ServiceProvider)
	a := r.Service(registry.ServiceNameApp).(app.App)

	srv := asynq.NewServer(
		cfg.AsynqConfig.RedisOpts(),
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: cfg.AsynqConfig.WorkerConcurrency,
			// Optionally specify multiple queues with different priority.
			Queues: cfg.AsynqConfig.Queues(),
		},
	)

	w := hsynq.NewWorker(srv, hexa.NewContextPropagator(sp.Logger(), sp.Translator()), hsynq.NewJsonTransformer())
	if err := jobs.RegisterJobs(w, sp, a); err != nil {
		return tracer.Trace(err)
	}

	registry.Register(registry.ServiceNameWorker, w)
	return nil
}