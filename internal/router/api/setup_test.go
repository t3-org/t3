package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/sr"
	"github.com/stretchr/testify/require"
	mockapp "space.org/space/internal/app/mock"
	"space.org/space/internal/registry"
	"space.org/space/internal/testbox"
)

var _any = gomock.Any()

// setup provides testReporter, mocked store and the app.
func setup(t *testing.T) (hexa.ServiceRegistry, *mockapp.MockApp) {
	r := sr.New()
	r = sr.NewMultiSearchRegistry(r, r, testbox.Global().Registry())
	r.Register(registry.ServiceNameTestReporter, t)

	require.NoError(t, registry.ProvideByNames(r, registry.ProviderNameMockApp))
	require.NoError(t, r.Boot())

	t.Cleanup(func() { require.NoError(t, r.Shutdown(context.Background())) })

	return r, r.Service(registry.ServiceNameApp).(*mockapp.MockApp)
}

func replyjson(code string, data string) string {
	return fmt.Sprintf(`{"code":"%s","data":%s}`, code, data)
}
