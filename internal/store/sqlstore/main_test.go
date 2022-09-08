package sqlstore_test

import (
	"os"
	"testing"
	"time"

	"github.com/kamva/gutil"
	huner "github.com/kamva/hexa-tuner"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/provider"
	"space.org/space/internal/testbox"
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
	tbox := testbox.New(provider.Providers(names))
	testbox.SetGlobal(tbox)

	if err = tbox.Setup(); err != nil {
		return
	}

	exitcode = t.Run()
	err = tracer.Trace(tbox.Teardown())
	return
}