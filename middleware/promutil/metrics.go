package promutil

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"strings"
)

func FindCollectorByName(collectors []MetricInfo, n string) MetricInfo {
	for _, c := range collectors {
		if c.Name == n {
			return c
		}
	}

	return MetricInfo{}
}

func FindCollectorById(collectors []MetricInfo, id string) MetricInfo {
	for _, c := range collectors {
		if c.Id == id {
			return c
		}
	}

	return MetricInfo{}
}

func SetMetricValueByName(collectors []MetricInfo, n string, v float64, labels prometheus.Labels) error {
	if c := FindCollectorByName(collectors, n); c.Type != "" {
		setMetricValue(c, v, labels)
	} else {
		err := errors.New("cannot find collector by name")
		log.Error().Err(err).Str("name", n).Send()
		return err
	}

	return nil
}

func SetMetricValueById(collectors []MetricInfo, id string, v float64, labels prometheus.Labels) error {
	if c := FindCollectorById(collectors, id); c.Type != "" {
		setMetricValue(c, v, labels)
	} else {
		err := errors.New("cannot find collector by id")
		log.Error().Err(err).Str("id", id).Send()
		return err
	}

	return nil
}

func setMetricValue(c MetricInfo, v float64, labels prometheus.Labels) {

	switch c.Type {
	case MetricTypeCounter:
		cnter := c.Collector.(*prometheus.CounterVec)
		cnter.With(labels).Add(v)
	case MetricTypeGauge:
		gauger := c.Collector.(*prometheus.GaugeVec)
		gauger.With(labels).Set(v)
	case MetricTypeHistogram:
		hist := c.Collector.(*prometheus.HistogramVec)
		hist.With(labels).Observe(v)
	}

}

func NewCollector(namespace string, subsystem string, opName string, metricConfig *MetricConfig) (prometheus.Collector, error) {

	var c prometheus.Collector
	switch metricConfig.Type {
	case MetricTypeCounter:
		c = NewCounter(namespace, subsystem, opName, metricConfig)
	case MetricTypeGauge:
		c = NewGauge(namespace, subsystem, opName, metricConfig)
	case MetricTypeHistogram:
		c = NewHistogram(namespace, subsystem, opName, metricConfig)
	default:
		return nil, errors.New("unknown metric type: " + metricConfig.Type)
	}

	if c == nil {
		return nil, errors.New("cannot instantiate metric: " + metricConfig.Name)
	}

	return c, nil
}

func NewCounter(namespace string, subsystem string, opName string, counterMetrics *MetricConfig) prometheus.Collector /* *prometheus.CounterVec */ {

	if counterMetrics.Type != MetricTypeCounter {
		log.Error().Str("type", counterMetrics.Type).Msg("type mismatch, not a counter")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("counter metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: metricSubsystem,
			Name:      counterMetrics.Name,
			Help:      counterMetrics.Help,
		},
		strings.Split(counterMetrics.Labels, ","))

	err := prometheus.Register(c)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", counterMetrics.Name).Msg("counter already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", counterMetrics.Name).Msg("counter error")
		}
	}

	return c
}

func NewGauge(namespace string, subsystem string, opName string, gaugeMetrics *MetricConfig) prometheus.Collector /* *prometheus.CounterVec */ {

	if gaugeMetrics.Type != MetricTypeGauge {
		log.Error().Str("type", gaugeMetrics.Type).Msg("type mismatch, not a gauge")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("gauge metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	c := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: metricSubsystem,
			Name:      gaugeMetrics.Name,
			Help:      gaugeMetrics.Help,
		},
		strings.Split(gaugeMetrics.Labels, ","))

	err := prometheus.Register(c)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", gaugeMetrics.Name).Msg("counter already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", gaugeMetrics.Name).Msg("gauge error")
		}
	}

	return c
}

func NewHistogram(namespace string, subsystem string, opName string, histogramMetrics *MetricConfig) prometheus.Collector {

	if histogramMetrics.Type != MetricTypeHistogram {
		log.Error().Str("type", histogramMetrics.Type).Msg("type mismatch, not a histogram")
		return nil
	}

	if namespace == "" || subsystem == "" || opName == "" {
		log.Error().Msg("histogram metric not configured, skipping creation")
		return nil
	}

	metricSubsystem := subsystem
	if strings.Contains(subsystem, "%s") {
		metricSubsystem = fmt.Sprintf(subsystem, opName)
	}

	var bck []float64
	switch t := histogramMetrics.Buckets.Type; t {
	case DefaultMetricsDurationBucketsTypeLinear:
		bck = prometheus.LinearBuckets(histogramMetrics.Buckets.Start, histogramMetrics.Buckets.WidthFactor, histogramMetrics.Buckets.Count)
	case DefaultMetricsDurationBucketsTypeExponential:
		bck = prometheus.ExponentialBuckets(histogramMetrics.Buckets.Start, histogramMetrics.Buckets.WidthFactor, histogramMetrics.Buckets.Count)
	case DefaultMetricsDurationBucketsTypeDefault:
		bck = prometheus.DefBuckets
	}

	h := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: metricSubsystem,
		Name:      histogramMetrics.Name,
		Help:      histogramMetrics.Help,
		Buckets:   bck,
	}, strings.Split(histogramMetrics.Labels, ","))

	err := prometheus.Register(h)
	if err != nil {
		if aregerr, ok := err.(prometheus.AlreadyRegisteredError); ok {
			log.Warn().Err(err).Str("name", histogramMetrics.Name).Msg("histogram already registered")
			return aregerr.ExistingCollector
		} else {
			log.Error().Err(err).Str("name", histogramMetrics.Name).Msg("histogram error")
		}
	}

	return h
}
