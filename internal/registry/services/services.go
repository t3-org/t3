package services

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

// Services is a simple facade that provides services using
// service registry. It provides the base services,
// other microservices,...
type Services interface {
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

type services struct {
	r hexa.ServiceRegistry
}

func (s *services) Config() hexa.Config {
	srv, _ := s.r.Service(registry.ServiceNameConfig).(*config.Config)
	return srv
}

func (s *services) Logger() hlog.Logger {
	srv, _ := s.r.Service(registry.ServiceNameLogger).(hlog.Logger)
	return srv
}

func (s *services) Translator() hexa.Translator {
	srv, _ := s.r.Service(registry.ServiceNameTranslator).(hexa.Translator)
	return srv
}

func (s *services) ProbeServer() probe.Server {
	srv, _ := s.r.Service(registry.ServiceNameProbeServer).(probe.Server)
	return srv
}

func (s *services) HealthReporter() hexa.HealthReporter {
	srv, _ := s.r.Service(registry.ServiceNameHealthReporter).(hexa.HealthReporter)
	return srv
}

func (s *services) Jobs() hjob.Jobs {
	srv, _ := s.r.Service(registry.ServiceNameJobs).(hjob.Jobs)
	return srv
}

func (s *services) OpenTelemetry() htel.OpenTelemetry {
	srv, _ := s.r.Service(registry.ServiceNameOpenTelemetry).(htel.OpenTelemetry)
	return srv
}

func (s *services) Redis() *hredis.Service {
	srv, _ := s.r.Service(registry.ServiceNameRedis).(*hredis.Service)
	return srv
}

func (s *services) CacheProvider() hcache.Provider {
	srv, _ := s.r.Service(registry.ServiceNameCacheProvider).(hcache.Provider)
	return srv
}

func (s *services) CronJobs() hjob.CronJobs {
	srv, _ := s.r.Service(registry.ServiceNameCron).(hjob.CronJobs)
	return srv
}

func (s *services) HttpServer() *hecho.EchoService {
	srv, _ := s.r.Service(registry.ServiceNameHttpServer).(*hecho.EchoService)
	return srv
}

func (s *services) Worker() hjob.Worker {
	srv, _ := s.r.Service(registry.ServiceNameWorker).(hjob.Worker)
	return srv
}

// New returns a Services facade.
func New(r hexa.ServiceRegistry) Services {
	return &services{r}
}

var _ Services = &services{}
