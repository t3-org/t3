package registry

import "github.com/kamva/gutil"

// What should we add to the service registry?
// - all services in the service provider.
// - store + app
// - the service that we run in the cmd function(e.g., grpc server, event listener or workflow worker)
const (
	ServiceNameConfig         = "config"
	ServiceNameIDGenerator    = "id_generator"
	ServiceNameLogger         = "logger"
	ServiceNameTranslator     = "translator"
	ServiceNameProbeServer    = "probe_server"
	ServiceNameHealthReporter = "health_reporter_service"
	ServiceNamePrometheus     = "prometheus"
	ServiceNameTracerProvider = "tracer_provider"
	ServiceNameMeterProvider  = "meter_provider"
	ServiceNameOpenTelemetry  = "open_telemetry"
	ServiceNameRedis          = "redis"
	ServiceNameMarkdown       = "markdown"
	//ServiceNameCacheProvider  = "cache_provider"

	ServiceNameHttpServer      = "http_server"
	ServiceNameJobs            = "jobs"
	ServiceNameWorker          = "worker"
	ServiceNameCron            = "cron"
	ServiceNameMatrix          = "matrix"
	ServiceNameChannels        = "channels"
	ServiceNameMatrixBotServer = "matrix_bot_server" // listens to the evens.

	// non service-provider's services:

	ServiceNameStore = "store"
	ServiceNameApp   = "app"

	// services to use in tests

	ServiceNameTempDB       = "temp_db"
	ServiceNameTestReporter = "test_reporter" // gomock.TestReporter
)

// We use all service names as the provider name, unless provide multiple providers for a
// service(e.g., mock of a service), in that situation we provide the provider name here.

const (
	ProviderNameMockStore = "mock_store"
	ProviderNameMockApp   = "mock_app"
)

func bootPriority() []string {
	return []string{
		ServiceNameConfig,
		ServiceNameIDGenerator,
		ServiceNameTempDB,
		ServiceNameLogger,
		ServiceNameTranslator,
		ServiceNameHealthReporter,
		ServiceNameProbeServer,
		ServiceNamePrometheus,
		ServiceNameTracerProvider,
		ServiceNameMeterProvider,
		ServiceNameOpenTelemetry,
		ServiceNameRedis,
		ServiceNameMarkdown,
		//ServiceNameCacheProvider,

		ServiceNameStore,
		ServiceNameMatrix,
		ServiceNameChannels,
		ServiceNameApp,

		ServiceNameHttpServer,
		ServiceNameJobs,
		ServiceNameWorker,
		ServiceNameCron,
		ServiceNameMatrixBotServer,
	}
}

func BaseServices(exclude ...string) []string {
	l := []string{
		ServiceNameConfig,
		ServiceNameIDGenerator,
		ServiceNameLogger,
		ServiceNameTranslator,
		ServiceNameProbeServer,
		ServiceNameHealthReporter,
		ServiceNamePrometheus,
		ServiceNameTracerProvider,
		ServiceNameMeterProvider,
		ServiceNameOpenTelemetry,
		ServiceNameRedis,
		ServiceNameMarkdown,

		//ServiceNameCacheProvider,
		ServiceNameStore,
		ServiceNameMatrix,
		ServiceNameChannels,
		ServiceNameApp,
	}

	if len(exclude) != 0 {
		return gutil.Sub(l, exclude)
	}
	return l
}

func TestHelperServices() []string {
	return []string{
		ServiceNameTempDB,
	}
}

func MinimalServices() []string {
	return []string{
		ServiceNameConfig,
		ServiceNameIDGenerator,
		ServiceNameLogger,
	}
}
