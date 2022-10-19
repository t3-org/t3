package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"space.org/space/internal/model"
	"space.org/space/internal/registry"
	"space.org/space/internal/registry/services"
	"space.org/space/internal/testbox"
)

func service() services.Services {
	return services.New(testbox.Global().Registry())
}

// setup provides testReporter, mocked store and the app.
func setup(t *testing.T) {
	s := testbox.Global().Registry().Service(registry.ServiceNameStore).(model.Store)
	require.NoError(t, s.TruncateAllTables(context.Background()))

}
