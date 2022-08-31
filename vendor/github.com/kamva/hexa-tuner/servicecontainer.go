package huner

import (
	"github.com/kamva/hexa"
	arranger "github.com/kamva/hexa-arranger"
	hevent "github.com/kamva/hexa-event"
	"github.com/kamva/hexa-job"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/hexa/htel"
	"github.com/kamva/hexa/probe"
)

type (
	// BaseServiceContainer is the base service container to use in each microservice.
	BaseServiceContainer interface {
		SetConfig(config hexa.Config)
		SetLogger(logger hlog.Logger)
		SetTranslator(translator hexa.Translator)
		SetProbeServer(server probe.Server)
		SetHealthReporter(reporter hexa.HealthReporter)
		SetJobs(jobs hjob.Jobs)
		SetEmitter(emitter hevent.Emitter)
		SetArranger(arranger arranger.Arranger)
		SetDLM(dlm hexa.DLM)
		SetOpenTelemetry(otlm htel.OpenTelemetry)

		Config() hexa.Config
		Logger() hlog.Logger
		Translator() hexa.Translator
		ProbeServer() probe.Server
		HealthReporter() hexa.HealthReporter
		Jobs() hjob.Jobs
		Emitter() hevent.Emitter
		Arranger() arranger.Arranger
		DLM() hexa.DLM
		OpenTelemetry() htel.OpenTelemetry
	}

	// baseServiceContainer contains all of services in one place to manage our services.
	baseServiceContainer struct {
		// must specify that should panic when user request
		// to get a nil service or just return nil value.
		must bool

		config         hexa.Config
		log            hlog.Logger
		translator     hexa.Translator
		probeServer    probe.Server
		healthReporter hexa.HealthReporter
		jobs           hjob.Jobs
		emitter        hevent.Emitter
		arranger       arranger.Arranger
		dlm            hexa.DLM
		otlm           htel.OpenTelemetry
	}
)

func (c *baseServiceContainer) SetConfig(config hexa.Config)             { c.config = config }
func (c *baseServiceContainer) SetLogger(logger hlog.Logger)             { c.log = logger }
func (c *baseServiceContainer) SetTranslator(translator hexa.Translator) { c.translator = translator }
func (c *baseServiceContainer) SetProbeServer(server probe.Server)       { c.probeServer = server }
func (c *baseServiceContainer) SetHealthReporter(r hexa.HealthReporter)  { c.healthReporter = r }
func (c *baseServiceContainer) SetJobs(jobs hjob.Jobs)                   { c.jobs = jobs }
func (c *baseServiceContainer) SetEmitter(emitter hevent.Emitter)        { c.emitter = emitter }
func (c *baseServiceContainer) SetArranger(arranger arranger.Arranger)   { c.arranger = arranger }
func (c *baseServiceContainer) SetDLM(dlm hexa.DLM)                      { c.dlm = dlm }
func (c *baseServiceContainer) SetOpenTelemetry(otlm htel.OpenTelemetry) { c.otlm = otlm }

func (c *baseServiceContainer) Config() hexa.Config                 { return c.config }
func (c *baseServiceContainer) Logger() hlog.Logger                 { return c.log }
func (c *baseServiceContainer) Translator() hexa.Translator         { return c.translator }
func (c *baseServiceContainer) ProbeServer() probe.Server           { return c.probeServer }
func (c *baseServiceContainer) HealthReporter() hexa.HealthReporter { return c.healthReporter }
func (c *baseServiceContainer) Jobs() hjob.Jobs                     { return c.jobs }
func (c *baseServiceContainer) Emitter() hevent.Emitter             { return c.emitter }
func (c *baseServiceContainer) Arranger() arranger.Arranger         { return c.arranger }
func (c *baseServiceContainer) DLM() hexa.DLM                       { return c.dlm }
func (c *baseServiceContainer) OpenTelemetry() htel.OpenTelemetry   { return c.otlm }

// NewBaseServiceContainer returns new instance of the BaseServiceContainer.
func NewBaseServiceContainer() BaseServiceContainer {
	return &baseServiceContainer{}
}
