package hsynq

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
)

type jobs struct {
	cli         *asynq.Client
	p           hexa.ContextPropagator
	transformer Transformer
}

type worker struct {
	s           *asynq.Server
	mux         *asynq.ServeMux
	p           hexa.ContextPropagator
	transformer Transformer
	closeCh     chan struct{}
}

func NewJobs(cli *asynq.Client, p hexa.ContextPropagator, t Transformer) hjob.Jobs {
	return &jobs{
		cli:         cli,
		p:           p,
		transformer: t,
	}
}

func NewWorker(s *asynq.Server, p hexa.ContextPropagator, t Transformer) hjob.Worker {
	return &worker{
		s:           s,
		mux:         asynq.NewServeMux(),
		p:           p,
		transformer: t,
		closeCh:     make(chan struct{}),
	}
}

func (j *jobs) Push(c context.Context, job *hjob.Job) error {
	ctxData, err := j.p.Inject(c)
	if err != nil {
		return tracer.Trace(err)
	}
	b, err := j.transformer.BytesFromJob(ctxData, job.Payload)
	if err != nil {
		return tracer.Trace(err)
	}
	_, err = j.cli.Enqueue(asynq.NewTask(job.Name, b), asynq.Queue(job.Queue), asynq.MaxRetry(job.Retry), asynq.Timeout(job.Timeout))
	return tracer.Trace(err)
}

func (j *jobs) Shutdown(c context.Context) error {
	return tracer.Trace(j.cli.Close())
}

func (w *worker) Register(name string, handlerFunc hjob.JobHandlerFunc) error {
	w.mux.HandleFunc(name, func(ctx context.Context, task *asynq.Task) error {
		headers, p, err := w.transformer.PayloadFromBytes(task.Payload())
		if err != nil {
			return tracer.Trace(err)
		}

		ctx, err = w.p.Extract(ctx, headers)
		if err != nil {
			return tracer.Trace(err)
		}
		err = tracer.Trace(handlerFunc(ctx, p))
		if err != nil {
			hexa.Logger(ctx).Error("error in handling queue task", hexa.ErrFields(err)...)
		}
		return err
	})

	return nil
}

func (w *worker) Run() error {
	// Currently we can not use the server's Run method,
	// because the server in the `Run` method listen to
	// signals to shutdown the server, but without any
	// mutex, so when we get a sigterm,... it try to
	// shutdown, while service registry listen to
	// signals too and calls to its shutdown method
	// again, we should wait for its fix that checks
	// if user is shutting down, do not shutdown
	// again to go not get any error.
	// Currently we just simply use the start method.
	// and wait for the closeCh to close.

	//return tracer.Trace(w.s.Run(w.mux))
	if err := w.s.Start(w.mux); err != nil {
		return tracer.Trace(err)
	}

	<-w.closeCh
	return nil
}

func (w *worker) Shutdown(ctx context.Context) error {
	w.s.Shutdown()
	close(w.closeCh)
	return nil
}

var _ hjob.Jobs = &jobs{}
var _ hexa.Shutdownable = &jobs{}
var _ hexa.Shutdownable = &worker{}
