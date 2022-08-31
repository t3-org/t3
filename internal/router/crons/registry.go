package crons

import (
	"github.com/kamva/gutil"
	hjob "github.com/kamva/hexa-job"
	"github.com/kamva/tracer"
	"space.org/space/internal/app"
	"space.org/space/internal/base"
	"space.org/space/internal/config"
)

func RegisterCronJobs(crons hjob.CronJobs, sp base.ServiceProvider, a app.App) error {
	res := newResources(sp, a)
	cfg := sp.Config().(*config.Config)
	_, _ = res, cfg
	err := gutil.AnyErr(
	//crons.Register(cfg.LDAP.SyncSchedule, hjob.NewCronJob("shield.user.sync_ldap_users"), res.SyncLdapUsers),

	// Register other jobs here.
	)

	return tracer.Trace(err)
}
