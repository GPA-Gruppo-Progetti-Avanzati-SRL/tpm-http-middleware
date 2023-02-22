package mwregistry_test

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

var cfg = []byte(`
  mw-handler-registry:
    gin-mw-metrics:
      namespace: leas_cab
      subsystem: tokens
    gin-mw-tracing:
      tags:
        - name: request.id
          type: header
          value: request-id
        - name: lra-id
          type: header
          value: Long-Running-Action
    gin-mw-error:
      status-code-policy: if-unlisted 
      status-code-policy-ranges:
        - from: 200
          to: 299
`)

type AppConfig struct {
	MwRegistry mwregistry.HandlerCatalogConfig `yaml:"mw-handler-registry" mapstructure:"mw-handler-registry" json:"mw-handler-registry"`
}

func TestError(t *testing.T) {

	appCfg := AppConfig{}
	err := yaml.Unmarshal(cfg, &appCfg)
	require.NoError(t, err)
}
