package mwmetrics

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-common/util/promutil"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

func init() {
	const semLogContext = "metrics-middleware::init"
	log.Info().Msg(semLogContext)
	mwregistry.RegisterHandlerFactory(MetricsHandlerId, NewPromHttpMetricsHandler)
}

type PromHttpMetricsHandler struct {
	config     *PromHttpMetricsHandlerConfig
	collectors promutil.Group
}

func MustNewPromHttpMetricsHandler(cfg interface{}) mws.MiddlewareHandler {

	const semLogContext = "metrics-handler::must-new"
	h, err := NewPromHttpMetricsHandler(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	return h
}

func NewPromHttpMetricsHandler(cfg interface{}) (mws.MiddlewareHandler, error) {
	const semLogContext = "metrics-handler::new"
	tcfg := DefaultPromHttpMetricsHandlerConfig

	if cfg != nil && !reflect.ValueOf(cfg).IsNil() {
		switch typedCfg := cfg.(type) {
		case mwregistry.HandlerCatalogConfig:
			err := mapstructure.Decode(typedCfg, &tcfg)
			if err != nil {
				return nil, err
			}
		case map[string]interface{}:
			err := mapstructure.Decode(typedCfg, &tcfg)
			if err != nil {
				return nil, err
			}
		case *PromHttpMetricsHandlerConfig:
			tcfg = *typedCfg
		default:
			log.Warn().Msg(semLogContext + " unmarshal issue for tracing handler config")
		}
	} else {
		log.Info().Str("mw-id", MetricsHandlerId).Msg(semLogContext + " config null...reverting to default values")
	}

	if tcfg.RefMetrics != nil {
		// Using global registry...
		log.Info().Interface("ref-metrics", tcfg.RefMetrics).Msg(semLogContext + " using externally defined metrics")
	} else {
		if tcfg.Namespace == "" || tcfg.Subsystem == "" {
			tcfg = DefaultMetricsConfig
		} else {
			if len(tcfg.Collectors) == 0 {
				tcfg.Collectors = DefaultMetricsConfig.Collectors
			}
		}

		tcfg.RefMetrics = &promutil.MetricsConfigReference{
			GId:         promutil.MetricsConfigReferenceLocalGroup,
			CounterId:   "requests",
			HistogramId: "request-duration",
		}
	}

	log.Info().Str("mw-id", MetricsHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	if tcfg.RefMetrics.IsLocal() {
		mregistry, err := promutil.InitGroup(promutil.MetricGroupConfig{Namespace: tcfg.Namespace, Subsystem: tcfg.Subsystem, Collectors: tcfg.Collectors})
		return &PromHttpMetricsHandler{config: &tcfg, collectors: mregistry}, err
		/*
			collectors := make([]promutil.Metric, 0)

			for _, mCfg := range tcfg.Collectors {
				if mc, err := promutil.NewCollector(tcfg.Namespace, tcfg.Subsystem, mCfg.Name, &mCfg); err != nil {
					log.Error().Err(err).Str("name", mCfg.Name).Msg("error creating metric")
				} else {
					collectors = append(collectors, promutil.Metric{Type: mCfg.Type, Id: mCfg.Id, Name: mCfg.Name, Collector: mc, Labels: mCfg.Labels})
				}
			}
		*/
	}

	return &PromHttpMetricsHandler{config: &tcfg}, nil
}

func (h *PromHttpMetricsHandler) GetKind() string {
	return MetricsHandlerKind
}

func (m *PromHttpMetricsHandler) HandleFunc() gin.HandlerFunc {

	const semLogContext = "metrics-handler::handle-func"
	return func(c *gin.Context) {

		g, _, err := m.config.RefMetrics.ResolveGroup(m.collectors)
		if err != nil {
			log.Error().Err(err).Msg(semLogContext + " disabling mw metrics")
			m.config.RefMetrics.GId = "-"
		}

		beginOfMiddleware := time.Now()

		var lbls prometheus.Labels

		if m.config.RefMetrics.IsHistogramEnabled() {
			defer func(begin time.Time) {
				err = g.SetMetricValueById(m.config.RefMetrics.HistogramId, time.Since(begin).Seconds(), lbls)
				if err != nil {
					log.Error().Err(err).Msg(semLogContext + " setting metrics")
				}
			}(beginOfMiddleware)
		}

		if nil != c {
			c.Next()
		}

		lbls = metricsLabels(c, c.Request.URL.String(), "500")
		if m.config.RefMetrics.IsCounterEnabled() {
			lbls[MetricStatusCodeLabelId] = fmt.Sprintf("%d", c.Writer.Status())
			_ = g.SetMetricValueById(m.config.RefMetrics.CounterId, 1, lbls)
		}
	}

}

const (
	MetricsCustomLabels     = MetricsHandlerId + "-custom-labels"
	MetricEndpointLabelId   = "endpoint"
	MetricStatusCodeLabelId = "status-code"
)

func metricsLabels(c *gin.Context, ep string, sc string) prometheus.Labels {

	const semLogContext = "metrics-handler::metrics-labels"

	metricsLabels := prometheus.Labels{
		MetricEndpointLabelId:   ep,
		MetricStatusCodeLabelId: sc,
	}

	if customLbls, ok := c.Get(MetricsCustomLabels); ok {
		if lbls, ok := customLbls.(prometheus.Labels); ok {
			for n, v := range lbls {
				metricsLabels[n] = v
			}
		} else {
			log.Warn().Str("custom-labels", MetricsCustomLabels).Msg(semLogContext + " found context key but not a prometheus.Labels")
		}
	}
	return metricsLabels
}
