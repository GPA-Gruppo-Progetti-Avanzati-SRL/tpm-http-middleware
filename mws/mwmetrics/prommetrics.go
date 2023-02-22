package mwmetrics

import (
	"fmt"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws/mwmetrics/promutil"
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
	collectors []promutil.MetricInfo
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

	log.Info().Interface("cfg", &tcfg).Msg(semLogContext)

	if tcfg.Namespace == "" || tcfg.Subsystem == "" {
		tcfg = DefaultMetricsConfig
	} else {
		if len(tcfg.Collectors) == 0 {
			tcfg.Collectors = DefaultMetricsConfig.Collectors
		}
	}

	log.Info().Str("mw-id", MetricsHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	collectors := make([]promutil.MetricInfo, 0)

	for _, mCfg := range tcfg.Collectors {
		if mc, err := promutil.NewCollector(tcfg.Namespace, tcfg.Subsystem, mCfg.Name, &mCfg); err != nil {
			log.Error().Err(err).Str("name", mCfg.Name).Msg("error creating metric")
		} else {
			collectors = append(collectors, promutil.MetricInfo{Type: mCfg.Type, Id: mCfg.Id, Name: mCfg.Name, Collector: mc, Labels: mCfg.Labels})
		}
	}

	return &PromHttpMetricsHandler{config: &tcfg, collectors: collectors}, nil
}

func (h *PromHttpMetricsHandler) GetKind() string {
	return MetricsHandlerKind
}

func (m *PromHttpMetricsHandler) HandleFunc() gin.HandlerFunc {

	return func(c *gin.Context) {

		beginOfMiddleware := time.Now()

		var sc = "500"
		ep := c.Request.URL.String()

		defer func(begin time.Time) {
			promutil.SetMetricValueById(m.collectors, "request_duration", time.Since(begin).Seconds(), prometheus.Labels{"endpoint": ep, "status_code": sc})
		}(beginOfMiddleware)

		if nil != c {
			c.Next()
		}

		sc = fmt.Sprintf("%d", c.Writer.Status())
		_ = promutil.SetMetricValueById(m.collectors, "requests", 1, prometheus.Labels{"endpoint": ep, "status_code": sc})
	}

}
