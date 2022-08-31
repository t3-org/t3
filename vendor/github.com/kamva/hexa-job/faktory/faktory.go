package hexafaktory

import (
	"context"
	"errors"
	"time"

	"github.com/contribsys/faktory/client"
	faktoryworker "github.com/contribsys/faktory_worker_go"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
)

type (
	// jobs is implementation of hexa Jobs using faktory
	jobs struct {
		p          *client.Pool
		propagator hexa.ContextPropagator
	}

	// worker is implementation of hexa worker using faktory.
	worker struct {
		w *faktoryworker.Manager
		p hexa.ContextPropagator
	}
)

func (j *jobs) prepare(c context.Context, job *hjob.Job) *client.Job {
	if job.Queue == "" {
		job.Queue = "default"
	}

	ctxData, _ := j.propagator.Inject(c)
	return &client.Job{
		Jid:       client.RandomJid(),
		Type:      job.Name,
		Queue:     job.Queue,
		Args:      []interface{}{ctxData, gutil.StructToMap(job.Payload)},
		CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Retry:     &job.Retry,

		// We don't using this custom data in any middleware, but just put it here :)
		Custom: bytesMapToInterfaceMap(ctxData),
	}
}

func (j *jobs) Push(ctx context.Context, job *hjob.Job) error {
	if job == nil || job.Name == "" {
		return tracer.Trace(errors.New("job is not valid (enter job name please)"))
	}

	return j.p.With(func(conn *client.Client) error {
		return conn.Push(j.prepare(ctx, job))
	})
}

func (w *worker) handler(jobName string, h hjob.JobHandlerFunc) faktoryworker.Perform {
	return func(ctx context.Context, args ...interface{}) error {
		ctxMap := args[0].(map[string]interface{})
		payload := newMapPayload(args[1].(map[string]interface{}))

		bytesMap, err := interfaceMapToBytesMap(ctxMap)
		if err != nil {
			return tracer.Trace(err)
		}

		c, err := w.p.Extract(ctx, bytesMap)
		if err != nil {
			return tracer.Trace(err)
		}

		return h(c, payload)
	}
}

func (w *worker) Register(name string, h hjob.JobHandlerFunc) error {
	w.w.Register(name, w.handler(name, h))
	return nil
}

func (w *worker) Run() error {
	// Run method does not return.
	w.w.Run()

	return nil
}

// NewFaktoryJobsDriver returns new instance of Jobs driver for the faktory
func NewFaktoryJobsDriver(p *client.Pool, propagator hexa.ContextPropagator) hjob.Jobs {
	return &jobs{p: p, propagator: propagator}
}

// NewFaktoryWorkerDriver returns new instance of hexa Worker driver for the faktory
func NewFaktoryWorkerDriver(w *faktoryworker.Manager, propagator hexa.ContextPropagator) hjob.Worker {
	return &worker{
		w: w,
		p: propagator,
	}
}

var _ hjob.Jobs = &jobs{}
var _ hjob.Worker = &worker{}
