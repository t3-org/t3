package registry

// What should we add to the service registry?
// - all services in the service provider.
// - store + app
// - the service that we run in the cmd function(e.g., grpc server, event listener or workflow worker)
const (
	ServiceNameConfig          = "config"
	ServiceNameServiceProvider = "service_provider"
	ServiceNameLogger          = "logger"
	ServiceNameTranslator      = "translator"
	ServiceNameProbeServer     = "probe_server"
	ServiceNameHealthReporter  = "health_reporter_service"
	ServiceNameTracerProvider  = "tracer_provider"
	ServiceNameMeterProvider   = "meter_provider"
	ServiceNameOpenTelemetry   = "open_telemetry"
	ServiceNamePrometheus      = "prometheus"
	ServiceNameRedis           = "redis"
	ServiceNameCacheProvider   = "cache_provider"

	ServiceNameHttpServer = "http_server"
	ServiceNameJobs       = "jobs"
	ServiceNameWorker     = "worker"
	ServiceNameCron       = "cron"

	// non service-provider's services:

	ServiceNameStore = "store"
	ServiceNameApp   = "app"
)

func bootPriority() []string {
	return []string{
		ServiceNameConfig,
		ServiceNameServiceProvider,
		ServiceNameLogger,
		ServiceNameTranslator,
		ServiceNameHealthReporter,
		ServiceNameProbeServer,
		ServiceNameTracerProvider,
		ServiceNameMeterProvider,
		ServiceNameOpenTelemetry,
		ServiceNamePrometheus,
		ServiceNameRedis,
		ServiceNameCacheProvider,
		ServiceNameHttpServer,
		ServiceNameJobs,
		ServiceNameWorker,
		ServiceNameCron,
		ServiceNameStore,
		ServiceNameApp,
	}
}

func BaseServices() []string {
	return []string{
		ServiceNameConfig,
		ServiceNameServiceProvider,
		ServiceNameLogger,
		ServiceNameTranslator,
		ServiceNameProbeServer,
		ServiceNameHealthReporter,
		ServiceNameTracerProvider,
		ServiceNameMeterProvider,
		ServiceNameOpenTelemetry,
		ServiceNamePrometheus,
		ServiceNameRedis,
		ServiceNameCacheProvider,
		ServiceNameStore,
		ServiceNameApp,
	}
}
