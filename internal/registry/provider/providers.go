package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gomodule/redigo/redis"
	"github.com/hibiken/asynq"
	"github.com/kamva/hexa"
	hecho "github.com/kamva/hexa-echo"
	"github.com/kamva/hexa-job/hsynq"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/hlog/logdriver"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
	"github.com/kamva/tracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sony/sonyflake"
	"go.mau.fi/util/dbutil"
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
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"t3.org/t3/internal/app"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/model"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
	"t3.org/t3/internal/router/api"
	"t3.org/t3/internal/router/jobs"
	"t3.org/t3/internal/service/channel"
	"t3.org/t3/internal/service/channel/matrixch"
	"t3.org/t3/internal/service/channel/matrixch/command"
	"t3.org/t3/internal/store"
	"t3.org/t3/internal/store/sqlstore"
	"t3.org/t3/pkg/md"
)

func init() {
	registry.AddProvider(registry.ServiceNameConfig, registry.ServiceNameConfig, ConfigProvider)
	registry.AddProvider(registry.ServiceNameIDGenerator, registry.ServiceNameIDGenerator, IDGeneratorProvider)
	registry.AddProvider(registry.ServiceNameTempDB, registry.ServiceNameTempDB, TmpDBProvider)
	registry.AddProvider(registry.ServiceNameLogger, registry.ServiceNameLogger, LoggerProvider)
	registry.AddProvider(registry.ServiceNameTranslator, registry.ServiceNameTranslator, TranslatorProvider)
	registry.AddProvider(registry.ServiceNameProbeServer, registry.ServiceNameProbeServer, ProbeServerProvider)
	registry.AddProvider(registry.ServiceNameHealthReporter, registry.ServiceNameHealthReporter, HealthReporterProvider)
	registry.AddProvider(registry.ServiceNameTracerProvider, registry.ServiceNameTracerProvider, TracerProvider)
	registry.AddProvider(registry.ServiceNamePrometheus, registry.ServiceNamePrometheus, PrometheusProvider)
	registry.AddProvider(registry.ServiceNameMeterProvider, registry.ServiceNameMeterProvider, MeterProvider)
	registry.AddProvider(registry.ServiceNameOpenTelemetry, registry.ServiceNameOpenTelemetry, OpenTelemetryProvider)
	registry.AddProvider(registry.ServiceNameRedis, registry.ServiceNameRedis, RedisProvider)
	registry.AddProvider(registry.ServiceNameMarkdown, registry.ServiceNameMarkdown, MarkdownProvider)
	registry.AddProvider(registry.ServiceNameHttpServer, registry.ServiceNameHttpServer, HttpServerProvider)
	registry.AddProvider(registry.ServiceNameJobs, registry.ServiceNameJobs, JobsProvider)
	registry.AddProvider(registry.ServiceNameWorker, registry.ServiceNameWorker, WorkerProvider)
	registry.AddProvider(registry.ServiceNameChannelHomes, registry.ServiceNameChannelHomes, ChannelHomesProvider)
	registry.AddProvider(registry.ServiceNameDispatcher, registry.ServiceNameDispatcher, DispatcherProvider)
	registry.AddProvider(registry.ServiceNameStore, registry.ServiceNameStore, StoreProvider)
	registry.AddProvider(registry.ServiceNameApp, registry.ServiceNameApp, AppProvider)
	registry.AddProvider(registry.ProviderNameMockStore, registry.ServiceNameStore, MockStoreProvider)
	registry.AddProvider(registry.ProviderNameMockApp, registry.ServiceNameApp, MockAppProvider)
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

func IDGeneratorProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	var settings sonyflake.Settings
	if cfg.MachineID != nil {
		settings.MachineID = func() (uint16, error) {
			return *cfg.MachineID, nil
		}
	}

	sf := sonyflake.NewSonyflake(settings)

	if sf == nil {
		return tracer.Trace(errors.New("can not create sonyflake"))
	}

	r.Register(registry.ServiceNameIDGenerator, sf)
	model.SetIDGenerator(sf)
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
	jobsInstance := hsynq.NewJobs(asynq.NewClient(cfg.AsynqConfig.RedisOpts()), propagator, hsynq.NewJsonTransformer())

	r.Register(registry.ServiceNameJobs, jobsInstance)
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

	r.Register(registry.ServiceNameTracerProvider, htel.NewTracerProvider(tp))
	return nil
}
func PrometheusProvider(r hexa.ServiceRegistry) error {
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
	probeServer := r.Service(registry.ServiceNameProbeServer).(probe.Server)
	probeServer.Register("prometheus_metrics", "/prometheus/metrics", exporter.ServeHTTP, "Prometheus metrics")

	// Register prometheus exporter as another service to just keep it:
	r.Register(registry.ServiceNamePrometheus, exporter)
	return nil
}

func MeterProvider(r hexa.ServiceRegistry) error {
	if !conf(r).OpenTelemetry.Metric.Enabled { // if it's disabled, ignore it.
		return nil
	}
	exporter := r.Service(registry.ServiceNamePrometheus).(*prometheus.Exporter)
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
	l := r.Service(registry.ServiceNameLogger).(hlog.Logger)

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

	pool := &redis.Pool{
		Dial:         dial,
		TestOnBorrow: testOnBorrow,
		MaxIdle:      5,
		IdleTimeout:  time.Second * 120,
	}

	ping := func(_ context.Context) error {
		_, err := pool.Get().Do("PING")
		return tracer.Trace(err)
	}

	r.RegisterByDescriptor(&hexa.Descriptor{
		Name:     registry.ServiceNameRedis,
		Instance: pool,
		Health:   hexa.NewPingHealth(l, registry.ServiceNameRedis, ping, nil),
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
	e.Use(hecho.CurrentUserBySub(nil)) // Set user as guest for all requests.

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
	api.New(e, a, api.NewMiddlewares(e)).RegisterRoutes()

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

func MarkdownProvider(r hexa.ServiceRegistry) error {
	// create Markdown parser extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock

	// create HTML renderer options
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}

	r.Register(registry.ServiceNameMarkdown, md.NewMarkdown(extensions, opts))
	return nil
}

func ChannelHomesProvider(r hexa.ServiceRegistry) error {
	homes := make(map[string]channel.Home)
	chConf, err := config.LoadChannelsConfig(conf(r).ChannelsConfigFilePath)
	if err != nil {
		return tracer.Trace(err)
	}
	for name, homeCfg := range chConf.ChannelHomes {
		cli, err := channelHomeProvider(r, homeCfg)
		if err != nil {
			return tracer.Trace(err)
		}
		homes[name] = cli
		r.Register(registry.ServiceNamePrefixChannelClient+name, cli)
	}

	r.Register(registry.ServiceNameChannelHomes, homes)
	return nil
}

func channelHomeProvider(r hexa.ServiceRegistry, rawCfg config.ChannelHome) (channel.Home, error) {
	switch rawCfg.Type {
	case config.ChannelHomeTypeMatrix: // Currently we use type just here, convert them to constant if needed.
		var cfg config.MatrixHomeConfig
		if err := rawCfg.Config.Decode(&cfg); err != nil {
			return nil, tracer.Trace(err)
		}

		cfg.ResolveEnvs(config.ChannelsEnvSuffix)

		return matrixProvider(r, cfg)
	}
	// Add other channel types to the switch case.
	return nil, tracer.Trace(fmt.Errorf("invalid channel type: %s", rawCfg.Type))
}

func matrixProvider(r hexa.ServiceRegistry, mcfg config.MatrixHomeConfig) (channel.Home, error) {
	cfg := conf(r)

	s := r.Service(registry.ServiceNameStore).(model.Store)
	db, err := dbutil.NewWithDB(s.DBLayer().(sqlstore.SqlStore).DB(), cfg.DB.Driver)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	cli, err := mautrix.NewClient(mcfg.HomeServerAddr, "", "")
	if err != nil {
		return nil, tracer.Trace(err)
	}

	cryptoHelper, err := cryptohelper.NewCryptoHelper(cli, []byte(mcfg.PickleKey), db)
	if err != nil {
		return nil, tracer.Trace(err)
	}

	cryptoHelper.LoginAs = &mautrix.ReqLogin{
		Type:                     mautrix.AuthTypePassword,
		InitialDeviceDisplayName: "T3",
		Identifier: mautrix.UserIdentifier{
			Type:    mautrix.IdentifierType(mcfg.IdentifierType),
			Medium:  mcfg.Medium,
			Address: mcfg.Address,
			User:    mcfg.Username,
		},
		Password: mcfg.Password,
	}

	// If you want to use multiple clients with the same DB, you
	// should set a distinct database account ID for each one.
	cryptoHelper.DBAccountID = mcfg.DBAccountID
	if err := cryptoHelper.Init(); err != nil {
		return nil, tracer.Trace(err)
	}

	cli.Crypto = cryptoHelper
	markdown := r.Service(registry.ServiceNameMarkdown).(*md.Markdown)

	hOpts := matrixch.HomeOpts{
		// We prefix the key with "_" to set the label as an internal
		// label on the ticket. see ticket labels.
		KeyPrefix:     config.InternalLabelKeyPrefix + "mx.", // e.g., t.mx.
		OkEmoji:       mcfg.OKEmoji,
		CommandPrefix: mcfg.CommandPrefix,
		Client:        cli,
		KVStore:       s.TicketLabel(),
		MarkDown:      markdown,
		UI:            cfg.UI,
	}

	// Provide the server
	router := matrixch.NewRouter(mcfg.CommandPrefix)
	bootFn := func() error {
		a := r.Service(registry.ServiceNameApp).(app.App)

		// initialize the event handlers.
		command.RegisterCommands(&hOpts, a, router)
		return nil
	}
	srv := matrixch.NewServer(cli, router, bootFn) // command handler server

	return matrixch.New(srv, &hOpts), nil
}

func DispatcherProvider(r hexa.ServiceRegistry) error {
	// create Markdown parser extensions
	homes := r.Service(registry.ServiceNameChannelHomes).(map[string]channel.Home)
	chConf, err := config.LoadChannelsConfig(conf(r).ChannelsConfigFilePath)
	if err != nil {
		return tracer.Trace(err)
	}

	// make channels
	channels := make(map[string]channel.Channel)
	for name, ch := range chConf.Channels {
		switch chConf.ChannelHomes[ch.Home].Type {
		case config.ChannelHomeTypeMatrix:
			var cfg matrixch.ChannelOptions
			if err := ch.Config.Decode(&cfg); err != nil {
				return tracer.Trace(fmt.Errorf("can not decode channel config, channel: %s: %w", name, err))
			}
			channels[name] = channel.Channel{
				Home:   ch.Home,
				Config: cfg,
			}
			// Put other cases here please.
		}
	}

	d := channel.NewDispatcher(homes, channels, chConf.Policies)
	r.Register(registry.ServiceNameDispatcher, d)
	return nil
}
