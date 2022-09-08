package provider

import (
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	"space.org/space/internal/config"
	"space.org/space/internal/registry"
	"space.org/space/internal/store/sqlstore"
)

func TmpDBProvider(r hexa.ServiceRegistry) error {
	cfg := conf(r)

	dbname := config.AppName + "_test_" + gutil.RandStringWithType(6, gutil.RandTypeLowercaseAlphaNum)
	if err := cfg.DB.SetName(dbname); err != nil {
		return tracer.Trace(err)
	}

	username, err := cfg.DB.Username()
	if err != nil {
		return tracer.Trace(err)
	}

	tempdb, err := sqlstore.NewTempDB(cfg.Test.DBRootConn, username, dbname)
	if err != nil {
		return tracer.Trace(err)
	}

	r.Register(registry.ServiceNameTempDB, tempdb)
	return nil
}
