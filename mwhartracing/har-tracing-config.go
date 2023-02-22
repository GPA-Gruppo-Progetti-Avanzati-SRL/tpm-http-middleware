package mwhartracing

const (
	HarTracingHandlerId   = "mw-har-tracing"
	HarTracingHandlerKind = "mw-kind-har-tracing"
)

type HarTracingHandlerConfig struct {
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
