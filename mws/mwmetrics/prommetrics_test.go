package mwmetrics_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwmetrics"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestPromHttpMetricsHandlerConfig(t *testing.T) {

	b, err := yaml.Marshal(mwmetrics.DefaultMetricsConfig)
	require.NoError(t, err)

	t.Log(string(b))
}
