package hexarobfig

import (
	"context"

	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"github.com/robfig/cron/v3"
)

type (
	// CronJobsOptions is the options that can provide on
	// create new cron job registerer to push cron jobs.
	CronJobsOptions struct {
		CtxGenerator ContextGenerator
		Cron         *cron.Cron
		Jobs         hjob.Jobs
		Worker       hjob.Worker
		Logger       hlog.Logger // optional
	}

	// ContextGenerator is a generator to generate new context to push as job's context.
	ContextGenerator func(ctx context.Context) context.Context

	// emptyPayload is just a empty payload as each job's payload.
	emptyPayload struct{}

	// cronJobPusher implements the hjob.CronJobs interface by pushing each cronJob into a job queue.
	cronJobPusher struct {
		ctxGenerator ContextGenerator
		logger       hlog.Logger
		cron         *cron.Cron
		jobs         hjob.Jobs
		worker       hjob.Worker
	}
)

func (c *cronJobPusher) Register(spec string, cJob *hjob.CronJob, handler hjob.CronJobHandlerFunc) error {
	if err := c.addCron(spec, cJob, handler); err != nil {
		return err
	}
	return c.registerJobHandler(cJob.Name, handler)
}

func (c *cronJobPusher) Run() (<-chan error, error) {
	done := make(chan error)
	go func() {
		c.cron.Run()
		close(done)
	}()
	return done, nil
}

func (c *cronJobPusher) Shutdown(_ context.Context) error {
	c.cron.Stop()
	return nil
}

// addCron sets the cron job to push new job on each call to cron-job handler
func (c *cronJobPusher) addCron(spec string, cJob *hjob.CronJob, handler hjob.CronJobHandlerFunc) error {
	_, err := c.cron.AddFunc(spec, func() {
		err := c.jobs.Push(c.ctxGenerator(context.Background()), c.job(cJob))
		if err != nil {
			c.reportFailedPush(cJob)
		}
	})

	return tracer.Trace(err)
}

// registerJobHandler sets the job handler.
func (c *cronJobPusher) registerJobHandler(jobName string, handler hjob.CronJobHandlerFunc) error {
	if c.worker == nil || handler == nil {
		c.logger.Info("worker or handler is nil, so this is a no-handler cron job, skip registering handler")
		return nil
	}
	return c.worker.Register(jobName, func(ctx context.Context, payload hjob.Payload) error {
		return handler(ctx)
	})
}

// job convert Cron job to a regular job.
func (c *cronJobPusher) job(job *hjob.CronJob) *hjob.Job {
	return &hjob.Job{
		Name:    job.Name,
		Queue:   job.Queue,
		Retry:   job.Retry,
		Payload: emptyPayload{},
	}
}

func (c *cronJobPusher) reportFailedPush(job *hjob.CronJob) {
	if c.logger != nil {
		c.logger.Error("failed to push cron job to the jobs service.",
			hlog.String("job_queue", job.Queue), hlog.String("job_name", job.Name))
	}
}

// NewCronJobPusher returns new instance of the Cron Jobs. It pushes each cron job into a job queue.
func NewCronJobPusher(options CronJobsOptions) hjob.CronJobs {
	return &cronJobPusher{
		ctxGenerator: options.CtxGenerator,
		logger:       options.Logger,
		cron:         options.Cron,
		jobs:         options.Jobs,
		worker:       options.Worker,
	}
}

// Assertion
var _ hjob.CronJobs = &cronJobPusher{}
var _ hexa.Runnable = &cronJobPusher{}
var _ hexa.Shutdownable = &cronJobPusher{}
