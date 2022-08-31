package base

import (
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/hibiken/asynq"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hcache "github.com/kamva/hexa-cache"
	"github.com/kamva/hexa-job/hsynq"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
	"github.com/kamva/tracer"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-redisstore"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/pkg/hredis"
)

type (
	// ServiceProvider is the app service provider that  just
	// provide base services,other microservices services,...
	// It does not provide our app services, those must provide
	// from entry-point of the app(probably main function) not here.
	ServiceProvider interface {
		huner.BaseServiceContainer
		Redis() *hredis.Service
		CacheProvider() hcache.Provider
	}

	// serviceProvider implements the ServiceProvider
	serviceProvider struct {
		huner.BaseServiceContainer
		redis *hredis.Service
		cp    hcache.Provider
	}
)

func (sp *serviceProvider) provideTracerProvider(cfg *config.Config) (trace.TracerProvider, error) {
	tcfg := cfg.OpenTelemetry.Tracing
	if tcfg.NoopTracer {
		return trace.NewNoopTracerProvider(), nil
	}
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(tcfg.JaegerAddr)))
	if err != nil {
		return nil, err
	}

	sampler := tracesdk.AlwaysSample()
	if !cfg.Debug && !tcfg.AlwaysSample {
		// We use the ParentBased(AlwaysSample)	sampler in production.
		sampler = tracesdk.ParentBased(tracesdk.AlwaysSample())
	}

	return tracesdk.NewTracerProvider(
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
	), nil
}

func (sp *serviceProvider) provideMeterProvider(cfg *config.Config, probeServer probe.Server) (metric.MeterProvider, error) {
	mcfg := cfg.OpenTelemetry.Metric

	if !mcfg.Enabled {
		return metric.NewNoopMeterProvider(), nil
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
		return nil, tracer.Trace(err)
	}

	// Register probe handler
	probeServer.Register("prometheus_metrics", "/prometheus/metrics", exporter.ServeHTTP, "Prometheus metrics")

	// Register prometheus exporter as another service to just keep it:
	registry.Register(registry.PrometheusService, exporter)
	return exporter.MeterProvider(), nil
}

func (sp *serviceProvider) provideProbeServer(cfg *config.Config, healthReporter hexa.HealthReporter) probe.Server {
	probeServer := probe.NewServer(&http.Server{Addr: cfg.ProbeServerAddress}, http.NewServeMux())

	probe.RegisterHealthHandlers(probeServer, healthReporter)
	// Register other probe server handlers here.

	return probeServer
}

func provideRedisPool(cfg *config.Config) *hredis.Service {
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

	return &hredis.Service{
		Pool: &redis.Pool{
			Dial:         dial,
			TestOnBorrow: testOnBorrow,
			MaxIdle:      5,
			IdleTimeout:  time.Second * 120,
		},
	}
}

// defaultBaseServices generate new instance of base pack with default implementation of services.
func (sp *serviceProvider) boot() error {
	cfg, err := config.New()
	if err != nil {
		return tracer.Trace(err)
	}
	config.SetDefaultConfig(cfg)

	logger := gutil.Must(logdriver.NewStackLoggerDriver(cfg.StackLoggerOptions())).(hlog.Logger)
	hlog.SetGlobalLogger(logger)
	translator := huner.NewTranslator(cfg.I18nPath(), cfg.TranslateOptions())
	hexatranslator.SetGlobal(translator)

	healthReporter := hexa.NewHealthReporter()
	probeServer := sp.provideProbeServer(cfg, healthReporter)
	propagator := hexa.NewContextPropagator(logger, translator)
	jobs := hsynq.NewJobs(asynq.NewClient(cfg.AsynqConfig.RedisOpts()), propagator, hsynq.NewJsonTransformer())

	tp, err := sp.provideTracerProvider(cfg)
	if err != nil {
		return tracer.Trace(err)
	}

	mp, err := sp.provideMeterProvider(cfg, probeServer)
	if err != nil {
		return tracer.Trace(err)
	}

	otlm := htel.NewOpenTelemetry(tp, mp)

	//--------------------------------
	// Initialize Base services
	//--------------------------------
	sp.SetConfig(cfg)
	sp.SetLogger(logger)
	sp.SetTranslator(translator)
	sp.SetProbeServer(probeServer)
	sp.SetHealthReporter(healthReporter)
	sp.SetOpenTelemetry(otlm)
	sp.SetJobs(jobs)

	//--------------------------------
	// Initialize Base Services
	//--------------------------------
	sp.redis = provideRedisPool(cfg)
	if cfg.Cache.Enabled {
		sp.cp = hcache.NewRedisCacheProvider(&hcache.RedisOptions{
			Prefix:      cfg.Cache.Redis.Prefix,
			Pool:        sp.Redis().Pool,
			Marshaler:   hcache.GobMarshal,
			Unmarshaler: hcache.GobUnmarshal,
			DefaultTTL:  cfg.Cache.Redis.TTL(),
		})
	}

	// Register services:
	registry.Register(registry.ConfigService, cfg)
	registry.Register(registry.LoggerService, logger)
	registry.Register(registry.TranslatorService, translator)
	registry.Register(registry.ProbeServerService, probeServer)
	registry.Register(registry.HealthReporterService, healthReporter)
	registry.Register(registry.TracerProviderService, tp)
	registry.Register(registry.MetricProviderService, mp)
	registry.Register(registry.OpenTelemetryService, otlm)
	registry.Register(registry.RedisService, sp.redis)
	registry.Register(registry.JobsService, jobs)
	registry.Register(registry.CacheProviderService, sp.cp)
	return nil
}

func (sp *serviceProvider) Redis() *hredis.Service {
	return sp.redis
}

func (sp *serviceProvider) CacheProvider() hcache.Provider {
	return sp.cp
}

// NewServiceProvider returns new instance of the ServiceProvider
func NewServiceProvider() (ServiceProvider, error) {
	sp := &serviceProvider{BaseServiceContainer: huner.NewBaseServiceContainer()}
	if err := sp.boot(); err != nil {
		return nil, tracer.Trace(err)
	}

	return sp, nil
}

// RateLimiter creates a new rate limiter.
func RateLimiter(cfg config.RateLimit, pool *redis.Pool) (limiter.Store, error) {
	rcfg := &redisstore.Config{
		Tokens:   cfg.Tokens,
		Interval: cfg.Interval(),
		Dial:     nil, // We use the pool, don't need to the dial fn.
	}
	return redisstore.NewWithPool(rcfg, pool)
}
