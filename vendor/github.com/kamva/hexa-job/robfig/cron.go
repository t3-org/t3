package hexarobfig

import (
	"context"

	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/robfig/cron/v3"
)

// cronJob implements hjob.CronJobs interface.
type cronJob struct {
	ctxGenerator ContextGenerator
	cron         *cron.Cron
}

func (c *cronJob) Register(spec string, cJob *hjob.CronJob, handler hjob.CronJobHandlerFunc) error {
	_, err := c.cron.AddFunc(spec, func() {
		err := handler(c.ctxGenerator(context.Background()))
		if err != nil {
			hlog.Error("failed cron job", hlog.ErrStack(err))
		}
	})

	return tracer.Trace(err)
}

func (c *cronJob) Run() error {
	c.cron.Run()
	return nil
}

func (c *cronJob) Shutdown(_ context.Context) error {
	c.cron.Stop()
	return nil
}

// New returns new instance of the Cron Jobs.
func New(cg ContextGenerator, cron *cron.Cron) hjob.CronJobs {
	return &cronJob{
		ctxGenerator: cg,
		cron:         cron,
	}
}

// Assertion
var _ hjob.CronJobs = &cronJob{}
var _ hexa.Runnable = &cronJob{}
var _ hexa.Shutdownable = &cronJob{}
