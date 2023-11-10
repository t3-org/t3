package provider

import (
	"github.com/golang/mock/gomock"
	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/tracer"
	mockapp "t3.org/t3/internal/app/mock"
	"t3.org/t3/internal/config"
	"t3.org/t3/internal/model"
	mockmodel "t3.org/t3/internal/model/mock"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/store/sqlstore"
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

func MockStoreProvider(r hexa.ServiceRegistry) error {
	s := mockmodel.NewStore(r.Service(registry.ServiceNameTestReporter).(gomock.TestReporter))
	r.Register(registry.ServiceNameStore, s)
	model.SetStore(s) // Set global DB store on the model package

	return nil
}

func MockAppProvider(r hexa.ServiceRegistry) error {
	ctl := gomock.NewController(r.Service(registry.ServiceNameTestReporter).(gomock.TestReporter))
	r.Register(registry.ServiceNameApp, mockapp.NewMockApp(ctl))
	return nil
}
