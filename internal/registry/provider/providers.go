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
	"space.org/space/internal/config"
	"space.org/space/internal/model"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/services"
	"space.org/space/internal/router/api"
	"space.org/space/internal/router/crons"
	"space.org/space/internal/router/jobs"
	"space.org/space/internal/store"
	cachestore "space.org/space/internal/store/cache"
	"space.org/space/internal/store/sqlstore"
	"space.org/space/pkg/hredis"
)

var providers = map[string]registry.Provider{
	registry.ServiceNameConfig:         ConfigProvider,
	registry.ServiceNameTempDB:         TmpDBProvider,
	registry.ServiceNameLogger:         LoggerProvider,
	registry.ServiceNameTranslator:     TranslatorProvider,
	registry.ServiceNameProbeServer:    ProbeServerProvider,
	registry.ServiceNameHealthReporter: HealthReporterProvider,
	registry.ServiceNameTracerProvider: TracerProvider,
	registry.ServiceNameMeterProvider:  MeterProvider,
	registry.ServiceNameOpenTelemetry:  OpenTelemetryProvider,
	registry.ServiceNameRedis:          RedisProvider,
	registry.ServiceNameCacheProvider:  CacheProvider,
	registry.ServiceNameHttpServer:     HttpServerProvider,
	registry.ServiceNameJobs:           JobsProvider,
	registry.ServiceNameWorker:         WorkerProvider,
	registry.ServiceNameCron:           CronProvider,
	registry.ServiceNameStore:          StoreProvider,
	registry.ServiceNameApp:            AppProvider,
}

// Providers returns provides with the specified names in the provided array.
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

func ConfigProvider(r hexa.ServiceRegistry) error {
	// Initialize configs:
	cfg, err := config.New()
	if err != nil {
		return tracer.Trace(err)
	}

	config.SetDefaultConfig(cfg)
	r.Register(registry.ServiceNameConfig, cfg)
	return nil
}

func LoggerProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)
	l, err := logdriver.NewStackLoggerDriver(cfg.StackLoggerOptions())
	if err != nil {
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameLogger, l)
	hlog.SetGlobalLogger(l)
	return nil
}

func TranslatorProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	translator := huner.NewTranslator(cfg.I18nPath(), cfg.TranslateOptions())
	hexatranslator.SetGlobal(translator)
	r.Register(registry.ServiceNameTranslator, translator)
	return nil
}

func HealthReporterProvider(r hexa.ServiceRegistry) error {
	r.Register(registry.ServiceNameHealthReporter, hexa.NewHealthReporter())
	return nil
}

func ProbeServerProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	reporter := r.Service(registry.ServiceNameHealthReporter).(hexa.HealthReporter)
	probeServer := probe.NewServer(&http.Server{Addr: cfg.ProbeServerAddress}, http.NewServeMux())

	probe.RegisterHealthHandlers(probeServer, reporter)
	// Register other probe server handlers here.

	r.Register(registry.ServiceNameProbeServer, probeServer)
	return nil
}

func JobsProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	l := r.Service(registry.ServiceNameLogger).(hlog.Logger)
	t := r.Service(registry.ServiceNameTranslator).(hexa.Translator)

	propagator := hexa.NewContextPropagator(l, t)
	jobs := hsynq.NewJobs(asynq.NewClient(cfg.AsynqConfig.RedisOpts()), propagator, hsynq.NewJsonTransformer())

	r.Register(registry.ServiceNameJobs, jobs)
	return nil
}

func TracerProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

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
			semconv.ServiceNameKey.String(config.AppName),
			attribute.String("environment", cfg.Environment),
			attribute.String("service_instance", cfg.InstanceName),
		)),
	)

	r.Register(registry.ServiceNameTracerProvider, tp)
	return nil
}

func MeterProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

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
			semconv.ServiceNameKey.String(config.AppName),
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

func OpenTelemetryProvider(r hexa.ServiceRegistry) error {
	tp := r.Service(registry.ServiceNameTracerProvider).(trace.TracerProvider)
	mp := r.Service(registry.ServiceNameMeterProvider).(metric.MeterProvider)
	r.Register(registry.ServiceNameOpenTelemetry, htel.NewOpenTelemetry(tp, mp))
	return nil
}

func RedisProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

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

func CacheProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)
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

func StoreProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)
	svcs := services.New(r)
	var s model.Store

	s, err := sqlstore.New(svcs.Logger(), cfg.DB)
	if err != nil {
		hlog.Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return tracer.Trace(err)
	}

	if svcs.CacheProvider() != nil { // Add the cache layer.
		s = cachestore.New(r, s)
	}

	if cfg.OpenTelemetry.Tracing.TraceDB {
		s = store.NewTracingLayerStore("sql", svcs.OpenTelemetry().TracerProvider(), s)
	}

	r.Register(registry.ServiceNameStore, s)

	// Set global DB store on the model package:
	model.SetStore(s)

	return nil
}

func AppProvider(r hexa.ServiceRegistry) error {
	s := r.Service(registry.ServiceNameStore).(model.Store)

	a, err := app.NewWithAllLayers(r, s)
	if err != nil {
		hlog.Error("error", hlog.ErrStack(tracer.Trace(err)), hlog.Err(err))
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameApp, a)
	return nil
}

func CronProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	s := services.New(r)
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
			BaseLogger:     s.Logger(),
			BaseTranslator: s.Translator(),
		})
	}

	cronJobs := hexarobfig.New(ctxGen, cron.New())

	if err := crons.RegisterCronJobs(cronJobs, r, cfg, a); err != nil {
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameCron, cronJobs)
	return nil
}

func tuneEcho(cfg *config.Config, r hexa.ServiceRegistry) *echo.Echo {
	s := services.New(r)
	metricsCfg := hecho.MetricsConfig{
		MeterProvider: s.OpenTelemetry().MeterProvider(),
	}

	tracingCfg := hecho.TracingConfig{
		Propagator:     propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}),
		TracerProvider: s.OpenTelemetry().TracerProvider(),
		ServerName:     config.AppName,
	}

	e := echo.New()

	// set limit options
	e.Server.ReadTimeout = time.Millisecond * time.Duration(cfg.RequestReadTimeoutMs)
	e.Server.ReadHeaderTimeout = time.Millisecond * time.Duration(cfg.RequestReadHeaderTimeoutMs)
	e.Server.WriteTimeout = time.Millisecond * time.Duration(cfg.ResponseWriteTimeoutMs)
	e.Server.IdleTimeout = time.Millisecond * time.Duration(cfg.ConnectionIdleTimeoutMs)
	e.Server.MaxHeaderBytes = cfg.MaxHeaderSizeKb << 10 // kb to bytes

	e.HideBanner = true
	e.Logger = hecho.HexaToEchoLogger(s.Logger(), cfg.EchoLogLevel)
	e.Debug = cfg.Debug
	e.HTTPErrorHandler = hecho.HTTPErrorHandler(s.Logger(), s.Translator(), e.Debug)

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
	e.Use(hecho.HexaContext(s.Logger(), s.Translator()))

	// SetContextLogger set the echo logger on each echo's context.
	e.Use(hecho.SetContextLogger(cfg.EchoLogLevel))

	// Add more data to each trace span:
	e.Use(hecho.TracingDataFromUserContext())

	return e
}

func HttpServerProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)
	a := r.Service(registry.ServiceNameApp).(app.App)

	e := tuneEcho(cfg, r)

	// Register Routes
	(&api.API{
		Echo:  e,
		API:   e.Group("api/v1"),
		App:   a,
		Guest: hecho.GuestMiddleware(),
		Auth:  hecho.AuthMiddleware(),
		Debug: hecho.DebugMiddleware(e),
	}).RegisterRoutes()

	echoService := &hecho.EchoService{Echo: e, Address: cfg.ListeningAddress()}
	r.Register(registry.ServiceNameHttpServer, echoService)
	return nil
}

func WorkerProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	l := r.Service(registry.ServiceNameLogger).(hlog.Logger)
	t := r.Service(registry.ServiceNameTranslator).(hexa.Translator)
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

	w := hsynq.NewWorker(srv, hexa.NewContextPropagator(l, t), hsynq.NewJsonTransformer())
	if err := jobs.RegisterJobs(w, r, a); err != nil {
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameWorker, w)
	return nil
}
