package crons

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/config"
)

func RegisterCronJobs(crons hjob.CronJobs, r hexa.ServiceRegistry, cfg *config.Config, a app.App) error {
	res := newResources(r, a)
	_, _ = res, cfg
	err := gutil.AnyErr(
	//crons.Register(cfg.LDAP.SyncSchedule, hjob.NewCronJob("space.user.sync_ldap_users"), res.SyncLdapUsers),

	// Register other jobs here.
	)

	return tracer.Trace(err)
}
