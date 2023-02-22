package mwtracing

const (
	TracingHandlerId   = "mw-tracing"
	TracingHandlerKind = "mw-kind-tracing"

	TracingHandlerSourceTypeHeader = "header"
)

type TracingTags struct {
	Name   string `yaml:"name,omitempty"  mapstructure:"name,omitempty" json:"name,omitempty"`
	Source string `yaml:"type,omitempty"  mapstructure:"type,omitempty" json:"type,omitempty"`
	Value  string `yaml:"value,omitempty"  mapstructure:"value,omitempty" json:"value,omitempty"`
}

type TracingHandlerConfig struct {
	Tags []TracingTags `yaml:"tags,omitempty"  mapstructure:"tags,omitempty" json:"tags,omitempty"`
}

var DefaultTracingHandlerConfig = TracingHandlerConfig{}

func (h *TracingHandlerConfig) GetKind() string {
	return TracingHandlerKind
}

type TracingHandlerConfigOption func(*TracingHandlerConfig)
type TracingHandlerConfigBuilder struct {
	opts []TracingHandlerConfigOption
}

func (cb *TracingHandlerConfigBuilder) Build() *TracingHandlerConfig {
	c := DefaultTracingHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
