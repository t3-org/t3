package app

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/sr"
	"github.com/stretchr/testify/require"
	mockmodel "space.org/space/internal/model/mock"
	"space.org/space/internal/registry"
	"space.org/space/internal/testbox"
)

var _any = gomock.Any()

// setup provides testReporter, mocked store and the app.
func setup(t *testing.T) (hexa.ServiceRegistry, *mockmodel.MockStore, App) {
	r := sr.New()
	r = sr.NewMultiSearchRegistry(r, r, testbox.Global().Registry())
	r.Register(registry.ServiceNameTestReporter, t)

	require.NoError(t, registry.ProvideByNames(r, registry.ProviderNameMockStore, registry.ServiceNameApp))
	require.NoError(t, r.Boot())

	t.Cleanup(func() { require.NoError(t, r.Shutdown(context.Background())) })

	return r, r.Service(registry.ServiceNameStore).(*mockmodel.MockStore), r.Service(registry.ServiceNameApp).(App)
}
