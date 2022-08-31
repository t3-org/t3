package huner

import (
	faktory "github.com/contribsys/faktory/client"
	faktoryworker "github.com/contribsys/faktory_worker_go"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa-job"
	"github.com/kamva/hexa-job/faktory"
	"github.com/kamva/tracer"
)

// NewFaktoryJobsDriver generate new faktory driver for hexa jobs.
func NewFaktoryJobsDriver(propagator hexa.ContextPropagator, poolSize int) (hjob.Jobs, error) {
	p, err := faktory.NewPool(poolSize)

	if err != nil {
		return nil, tracer.Trace(err)
	}

	return hexafaktory.NewFaktoryJobsDriver(p, propagator), nil
}

// NewFaktoryWorkersDriver generate new faktory driver for the hexa worker.
func NewFaktoryWorkersDriver(p hexa.ContextPropagator, concurrency int) (hjob.Worker, error) {
	mgr := faktoryworker.NewManager()
	mgr.Concurrency = concurrency
	worker := hexafaktory.NewFaktoryWorkerDriver(mgr, p)
	return worker, nil
}
