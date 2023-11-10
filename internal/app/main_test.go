package app_test

import (
	"os"
	"testing"
	"time"

	_ "t3.org/t3/internal/registry/provider"

	"github.com/kamva/gutil"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/testbox"
)

// TODO: this function is repeated in other packages like sqlstore/main-test.go, move to somewhere like testbox package.
func TestMain(t *testing.M) {
	// Enable Test environment.
	gutil.PanicErr(os.Setenv(huner.EnvironmentKey(huner.EnvKeysPrefix()), huner.EnvironmentTest))

	code, err := testMain(t)
	if err != nil {
		hlog.Error("error on running tests", hlog.ErrStack(err))
		time.Sleep(time.Second)
		os.Exit(1)
	}
	os.Exit(code)
}

func testMain(t *testing.M) (exitcode int, err error) {
	baseServices := registry.BaseServices(registry.ServiceNameStore, registry.ServiceNameApp)
	providers, err := registry.Providers(baseServices...)
	if err != nil {
		return 0, tracer.Trace(err)
	}

	tbox := testbox.New(providers)
	testbox.SetGlobal(tbox)

	if err = tbox.Setup(); err != nil {
		return
	}

	exitcode = t.Run()
	err = tracer.Trace(tbox.Teardown())
	return
}
