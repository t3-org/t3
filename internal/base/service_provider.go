package base

import (
	"github.com/kamva/hexa"
	hcache "github.com/kamva/hexa-cache"
	hecho "github.com/kamva/hexa-echo"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/pkg/hredis"
)

// ServiceProvider is the app service provider that  just
// provide base services,other microservices services,...
// It does not provide our app services, those must provide
// from entry-point of the app(probably main function) not here.
type ServiceProvider interface {
	Config() hexa.Config
	Logger() hlog.Logger
	Translator() hexa.Translator
	ProbeServer() probe.Server
	HealthReporter() hexa.HealthReporter
	Jobs() hjob.Jobs
	OpenTelemetry() htel.OpenTelemetry
	Redis() *hredis.Service
	CacheProvider() hcache.Provider
	CronJobs() hjob.CronJobs
	HttpServer() *hecho.EchoService
	Worker() hjob.Worker
}

type serviceProvider struct {
	r hexa.ServiceRegistry
}

func (p *serviceProvider) Config() hexa.Config {
	srv, _ := p.r.Service(registry.ServiceNameConfig).(*config.Config)
	return srv
}

func (p *serviceProvider) Logger() hlog.Logger {
	srv, _ := p.r.Service(registry.ServiceNameLogger).(hlog.Logger)
	return srv
}

func (p *serviceProvider) Translator() hexa.Translator {
	srv, _ := p.r.Service(registry.ServiceNameTranslator).(hexa.Translator)
	return srv
}

func (p *serviceProvider) ProbeServer() probe.Server {
	srv, _ := p.r.Service(registry.ServiceNameProbeServer).(probe.Server)
	return srv
}

func (p *serviceProvider) HealthReporter() hexa.HealthReporter {
	srv, _ := p.r.Service(registry.ServiceNameHealthReporter).(hexa.HealthReporter)
	return srv
}

func (p *serviceProvider) Jobs() hjob.Jobs {
	srv, _ := p.r.Service(registry.ServiceNameJobs).(hjob.Jobs)
	return srv
}

func (p *serviceProvider) OpenTelemetry() htel.OpenTelemetry {
	srv, _ := p.r.Service(registry.ServiceNameOpenTelemetry).(htel.OpenTelemetry)
	return srv
}

func (p *serviceProvider) Redis() *hredis.Service {
	srv, _ := p.r.Service(registry.ServiceNameRedis).(*hredis.Service)
	return srv
}

func (p *serviceProvider) CacheProvider() hcache.Provider {
	srv, _ := p.r.Service(registry.ServiceNameCacheProvider).(hcache.Provider)
	return srv
}

func (p *serviceProvider) CronJobs() hjob.CronJobs {
	srv, _ := p.r.Service(registry.ServiceNameCron).(hjob.CronJobs)
	return srv
}

func (p *serviceProvider) HttpServer() *hecho.EchoService {
	srv, _ := p.r.Service(registry.ServiceNameHttpServer).(*hecho.EchoService)
	return srv
}

func (p *serviceProvider) Worker() hjob.Worker {
	srv, _ := p.r.Service(registry.ServiceNameWorker).(hjob.Worker)
	return srv
}

// NewServiceProvider returns new instance of the ServiceProvider
func NewServiceProvider(r hexa.ServiceRegistry) ServiceProvider {
	return &serviceProvider{r}
}

var _ ServiceProvider = &serviceProvider{}
