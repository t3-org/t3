package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatrixHomeConfig_ResolveEnvs(t *testing.T) {
	c := MatrixHomeConfig{
		HomeServerAddr: "__ENV__HOME_ADDR",
		Address:        "__ENV__ADDRESS",
		Username:       "__ENV__USERNAME",
		Password:       "pass123",
	}

	require.NoError(t, os.Setenv("__ENV__HOME_ADDR", "home123"))
	require.NoError(t, os.Setenv("__ENV__ADDRESS", "email123"))
	require.NoError(t, os.Setenv("__ENV__USERNAME", "username123"))

	c.ResolveEnvs("__ENV__")

	require.Equal(t, "home123", c.HomeServerAddr)
	require.Equal(t, "email123", c.Address)
	require.Equal(t, "username123", c.Username)
	require.Equal(t, "pass123", c.Password)

	// Set just one field
	c.Password = "__ENV__PASSWORD"
	require.NoError(t, os.Setenv("__ENV__PASSWORD", "pass456"))

	// Make sure it doesn't resolve any un-prefixed value.
	c.ResolveEnvs("__ENV__")
	require.Equal(t, "home123", c.HomeServerAddr)
	require.Equal(t, "email123", c.Address)
	require.Equal(t, "username123", c.Username)
	require.Equal(t, "pass456", c.Password)
}
