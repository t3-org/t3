package registry

// What should we add to the service registry?
// - all services in the service provider.
// - store + app
// - the service that we run in the cmd function(e.g., grpc server, event listener or workflow worker)
const (
	ConfigService         = "config"
	LoggerService         = "logger"
	TranslatorService     = "translator"
	ProbeServerService    = "probe_server"
	HealthReporterService = "health_reporter_service"
	TracerProviderService = "tracer_provider"
	MetricProviderService = "metric_provider"
	OpenTelemetryService  = "open_telemetry"
	PrometheusService     = "prometheus"
	RedisService          = "redis"
	DLMService            = "dlm"

	HttpServerService    = "http_server"
	JobsService          = "jobs"
	WorkerService        = "worker"
	CronService          = "cron"
	CacheProviderService = "cache_provider"

	// non service-provider's services:

	StoreService = "store"
	AppService   = "app"
)
