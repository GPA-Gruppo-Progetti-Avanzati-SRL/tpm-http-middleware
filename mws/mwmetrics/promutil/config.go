package promutil_old

import (
	"github.com/prometheus/client_golang/prometheus"
)

const DefaultMetricsDurationBucketsTypeLinear = "linear"
const DefaultMetricsDurationBucketsTypeExponential = "exponential"
const DefaultMetricsDurationBucketsTypeDefault = "default"

const DefaultMetricsDurationBucketsStart = 0.5
const DefaultMetricsDurationBucketsWidthFormat = 0.5
const DefaultMetricsDurationBucketsCount = 10

const MetricTypeCounter = "counter"
const MetricTypeGauge = "gauge"
const MetricTypeHistogram = "histogram"

//type MetricsCounterConfig struct {
//	Name   string
//	Help   string
//	Labels string
//}
//
//type MetricsGaugeConfig struct {
//	Name   string
//	Help   string
//	Labels string
//}

type MetricInfo struct {
	Id        string
	Type      string
	Name      string
	Collector prometheus.Collector
	Labels    string
}

type MetricsConfig struct {
	Namespace  string
	Subsystem  string
	Collectors []MetricConfig
}

type MetricConfig struct {
	Id      string                `yaml:"id"  mapstructure:"id"  json:"id"`
	Name    string                `yaml:"name"  mapstructure:"name"  json:"name"`
	Help    string                `yaml:"help"  mapstructure:"help"  json:"help"`
	Labels  string                `yaml:"labels"  mapstructure:"labels"  json:"labels"`
	Type    string                `yaml:"type"  mapstructure:"type"  json:"type"`
	Buckets HistogramBucketConfig `yaml:"buckets,omitempty"  mapstructure:"buckets,omitempty"  json:"buckets,omitempty"`
}

/*type MetricsHistogramConfig struct {
	Name    string
	Help    string
	Labels  string
	Buckets HistogramBucketConfig
}
*/

type HistogramBucketConfig struct {
	Type        string  `yaml:"type,omitempty"  mapstructure:"type,omitempty"  json:"type,omitempty"`
	Start       float64 `yaml:"start,omitempty"  mapstructure:"start,omitempty"  json:"start,omitempty"`
	WidthFactor float64 `yaml:"width-factor,omitempty"  mapstructure:"width-factor,omitempty"  json:"width-factor,omitempty"`
	Count       int     `yaml:"count,omitempty"  mapstructure:"count,omitempty"  json:"count,omitempty"`
}
