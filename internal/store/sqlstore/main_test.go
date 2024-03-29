package sqlstore_test

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
	names := append(registry.MinimalServices(), registry.TestHelperServices()...)
	names = append(names, registry.ServiceNameStore) // Add store.
	providers, err := registry.Providers(names...)
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
