package jobs

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
)

func RegisterJobs(w hjob.Worker, r hexa.ServiceRegistry, a app.App) error {
	res := newResources(r, a)
	_ = w
	_ = res
	err := gutil.AnyErr(
	//w.Register(model.JobNameSyncLdapUsers, res.SyncLdapUsers),

	// Register other jobs here.
	)

	return tracer.Trace(err)
}
