package arranger

import (
	"context"

	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"go.temporal.io/sdk/worker"
)

// Worker is just a wrapper for worker.Worker implement hexa service.
type Worker interface {
	// Worker returns the temporal worker.
	// We could not embed temporal worker because of cnflict
	// with the Run() method.
	Worker() worker.Worker
	hexa.Runnable
	hexa.Shutdownable
}

type workerImpl struct {
	w worker.Worker
}

func NewWorker(w worker.Worker) Worker {
	return &workerImpl{w: w}
}

func (w *workerImpl) Worker() worker.Worker {
	return w.w
}

func (w *workerImpl) Run() error {
	return tracer.Trace(w.w.Run(nil))
}

func (w *workerImpl) Shutdown(c context.Context) error {
	w.w.Stop()
	return nil
}

var _ hexa.Runnable = &workerImpl{}
var _ hexa.Shutdownable = &workerImpl{}
