package services

import (
	"github.com/gomodule/redigo/redis"
	"github.com/kamva/hexa"
	hecho "github.com/kamva/hexa-echo"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/service/channel"
)

// Services is a simple facade that provides services using
// service registry. It provides the base services,
// other microservices,...
type Services interface {
	Config() *config.Config
	Logger() hlog.Logger
	Translator() hexa.Translator
	ProbeServer() probe.Server
	HealthReporter() hexa.HealthReporter
	Jobs() hjob.Jobs
	OpenTelemetry() htel.OpenTelemetry
	Redis() *redis.Pool
	CronJobs() hjob.CronJobs
	HttpServer() *hecho.EchoService
	Worker() hjob.Worker
	Matrix() *channel.MatrixChannel
	Channels() map[string]channel.Channel
}

type services struct {
	r hexa.ServiceRegistry
}

func (s *services) Config() *config.Config {
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

func (s *services) Redis() *redis.Pool {
	srv, _ := s.r.Service(registry.ServiceNameRedis).(*redis.Pool)
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

func (s *services) Matrix() *channel.MatrixChannel {
	return s.r.Service(registry.ServiceNameMatrix).(*channel.MatrixChannel)
}

func (s *services) Channels() map[string]channel.Channel {
	return s.r.Service(registry.ServiceNameChannels).(map[string]channel.Channel)
}

// New returns a Services facade.
func New(r hexa.ServiceRegistry) Services {
	return &services{r}
}

var _ Services = &services{}
