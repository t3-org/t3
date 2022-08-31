package hjob

import (
	"context"

	"github.com/kamva/hexa"
)

type (
	// CronJobHandlerFunc is the handler of each cron job in the worker.
	CronJobHandlerFunc func(ctx context.Context) error

	// CronJob is a new instance of cron job that should run in schedules time.
	CronJob struct {
		Name  string // required
		Queue string
		// Retry specify retry counts of the job.
		// 0: means that throw job away (and dont push to dead queue) on first fail.
		// -1: means that push job to the dead queue on first fail.
		Retry int
	}

	// CronJobs get your cron-jobs specs
	CronJobs interface {
		// Register handler for new cron job
		Register(spec string, cJob *CronJob, handler CronJobHandlerFunc) error
		hexa.Runnable
	}
)

// NewCronJob returns new cron job instance
func NewCronJob(name string) *CronJob {
	return NewCronJobWithQueue(name, "default")
}

// NewCronJobWithQueue returns new cron job instance
func NewCronJobWithQueue(name string, queue string) *CronJob {
	return &CronJob{
		Name:  name,
		Queue: queue,
		Retry: 4,
	}
}
