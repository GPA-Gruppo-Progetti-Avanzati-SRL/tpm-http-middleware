package middleware_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/middleware"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestPromHttpMetricsHandlerConfig(t *testing.T) {

	b, err := yaml.Marshal(middleware.DefaultMetricsConfig)
	require.NoError(t, err)

	t.Log(string(b))
}
