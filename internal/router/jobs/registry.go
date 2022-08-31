package jobs

import (
	"github.com/kamva/gutil"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
)

func RegisterJobs(w hjob.Worker, sp base.ServiceProvider, a app.App) error {
	res := newResources(sp, a)
	_ = res
	err := gutil.AnyErr(
		//w.Register(model.JobNameSyncLdapUsers, res.SyncLdapUsers),

		// Register other jobs here.
	)

	return tracer.Trace(err)
}
