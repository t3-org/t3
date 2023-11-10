package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"t3.org/t3/internal/model"
	"t3.org/t3/internal/registry"
	"t3.org/t3/internal/registry/services"
	"t3.org/t3/internal/testbox"
)

func service() services.Services {
	return services.New(testbox.Global().Registry())
}

// setup provides testReporter, mocked store and the app.
func setup(t *testing.T) {
	s := testbox.Global().Registry().Service(registry.ServiceNameStore).(model.Store)
	require.NoError(t, s.TruncateAllTables(context.Background()))

}
