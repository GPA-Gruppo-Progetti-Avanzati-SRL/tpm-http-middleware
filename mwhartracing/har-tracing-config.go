package mwhartracing

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"

const (
	HarTracingHandlerId   = "mw-har-tracing"
	HarTracingHandlerKind = "mw-kind-har-tracing"
)

type HarTracingHandlerConfig struct {
	RefMetrics *promutil.MetricsConfigReference `yaml:"ref-metrics"  mapstructure:"ref-metrics"  json:"ref-metrics"`
}

var DefaultTracingHandlerConfig = HarTracingHandlerConfig{}

func (h *HarTracingHandlerConfig) GetKind() string {
	return HarTracingHandlerKind
}

type TracingHandlerConfigOption func(*HarTracingHandlerConfig)
type TracingHandlerConfigBuilder struct {
	opts []TracingHandlerConfigOption
}

func (cb *TracingHandlerConfigBuilder) Build() *HarTracingHandlerConfig {
	c := DefaultTracingHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
