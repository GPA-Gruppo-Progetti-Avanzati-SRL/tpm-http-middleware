package mwmetrics

import "github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"

const (
	MetricsHandlerId   = "mw-metrics"
	MetricsHandlerKind = "mw-kind-metrics"
)

var DefaultMetricsConfig = PromHttpMetricsHandlerConfig{
	Namespace: "tpm",
	Subsystem: "http_server",
	Collectors: []promutil.MetricConfig{
		{
			Id:   "requests",
			Name: "requests",
			Help: "numero richieste",
			Labels: promutil.MetricLabelsConfig{
				{
					Id:           "endpoint",
					Name:         "endpoint",
					DefaultValue: "N/A",
				},
				{
					Id:           "status-code",
					Name:         "status_code",
					DefaultValue: "-1",
				},
			},
			Type: promutil.MetricTypeCounter,
		},
		{
			Id:   "request_duration",
			Name: "request_duration",
			Help: "durata lavorazione richiesta",
			Labels: promutil.MetricLabelsConfig{
				{
					Id:           "endpoint",
					Name:         "endpoint",
					DefaultValue: "N/A",
				},
				{
					Id:           "status-code",
					Name:         "status_code",
					DefaultValue: "-1",
				},
			},
			Type: promutil.MetricTypeHistogram,
			Buckets: promutil.HistogramBucketConfig{
				Type:        "linear",
				Start:       promutil.DefaultMetricsDurationBucketsStart,
				WidthFactor: promutil.DefaultMetricsDurationBucketsWidthFormat,
				Count:       promutil.DefaultMetricsDurationBucketsCount,
			},
		},
	},
}

/*
 * ErrorHandlerConfig
 */

type PromHttpMetricsHandlerConfig struct {
	RefMetrics *promutil.MetricsConfigReference `yaml:"ref-metrics"  mapstructure:"ref-metrics"  json:"ref-metrics"`
	Namespace  string                           `yaml:"namespace"  mapstructure:"namespace"  json:"namespace"`
	Subsystem  string                           `yaml:"subsystem"  mapstructure:"subsystem"  json:"subsystem"`
	Collectors []promutil.MetricConfig          `yaml:"metrics"  mapstructure:"metrics"  json:"metrics"`
}

var DefaultPromHttpMetricsHandlerConfig = PromHttpMetricsHandlerConfig{}

func (h *PromHttpMetricsHandlerConfig) GetKind() string {
	return MetricsHandlerKind
}

type PromHttpMetricsHandlerOption func(*PromHttpMetricsHandlerConfig)
type PromHttpMetricsHandlerConfigBuilder struct {
	opts []PromHttpMetricsHandlerOption
}

/*
func (cb *PromHttpMetricsHandlerConfigBuilder) WithEndpoint(endpoint string) *PromHttpMetricsHandlerConfigBuilder {

	handlerFactoryMap := func(c *PromHttpMetricsHandlerConfig) {
		c.Endpoint = endpoint
	}

	cb.opts = append(cb.opts, handlerFactoryMap)
	return cb
}
*/

func (cb *PromHttpMetricsHandlerConfigBuilder) Build() *PromHttpMetricsHandlerConfig {
	c := DefaultPromHttpMetricsHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
