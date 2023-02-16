package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HandlerCatalogConfig map[string]HandlerConfig
type HandlerConfig map[string]interface{}

/*
struct {
	ErrCfg     *ErrorHandlerConfig           `yaml:"gin-mw-error" mapstructure:"gin-mw-error" json:"gin-mw-error"`
	MetricsCfg *PromHttpMetricsHandlerConfig `yaml:"gin-mw-metrics" mapstructure:"gin-mw-metrics" json:"gin-mw-metrics"`
	TraceCfg   *TracingHandlerConfig         `yaml:"gin-mw-tracing" mapstructure:"gin-mw-tracing" json:"gin-mw-tracing"`
}
*/

type HandlerFactory func(interface{}) (MiddlewareHandler, error)

var handlerFactoryMap = map[string]HandlerFactory{
	ErrorHandlerId:   NewErrorHandler,
	TracingHandlerId: NewTracingHandler,
	MetricsHandlerId: NewPromHttpMetricsHandler,
}

func RegisterHandlerFactory(handlerId string, hf HandlerFactory) {
	const semLogContext = "middleware:register-handler-factory"
	if _, ok := handlerFactoryMap[handlerId]; ok {
		log.Warn().Str("mw-id", handlerId).Msg(semLogContext + " handler factory already registered")
		return
	}

	handlerFactoryMap[handlerId] = hf
}

type HandlerRegistry map[string]gin.HandlerFunc

var registry HandlerRegistry = make(map[string]gin.HandlerFunc)

func InitializeHandlerRegistry(registryConfig HandlerCatalogConfig, mwInUse []string) error {

	const semLogContext = "middleware:registry-initialization"

	for _, mw := range mwInUse {

		factory, ok := handlerFactoryMap[mw]
		if !ok {
			log.Error().Str("mw-id", mw).Msg(semLogContext + " cannot find middleware in catalog")
			continue
		}

		log.Info().Str("mw-id", mw).Msg(semLogContext + " initializing handler")
		cfg := registryConfig[mw]
		r, err := factory(cfg)
		if err != nil {
			log.Error().Err(err).Str("mw-id", mw).Msg(semLogContext + " initialization handler failure")
			continue
		}

		registry[mw] = r.HandleFunc()

		/*
			switch mw {
			case ErrorHandlerId:
				registry[ErrorHandlerId] = NewErrorHandler(registryConfig.ErrCfg).HandleFunc()
			case TracingHandlerId:
				registry[TracingHandlerId] = NewTracingHandler(registryConfig.TraceCfg).HandleFunc()
			case MetricsHandlerId:
				registry[MetricsHandlerId] = NewPromHttpMetricsHandler(registryConfig.MetricsCfg).HandleFunc()
			}
		*/
	}

	/*
		for n, i := range registryConfig {
			if hanlderFactory, ok := handlerFactoryMap[n]; ok {
				registry[n] = hanlderFactory(i).HandleFunc()
			} else {
				err := errors.New("cannot find factory for middleware handler of id: " + n)
				log.Error().Err(err).Send()
				return err
			}
		}
	*/

	return nil
}

func GetHandlerFunc(name string) gin.HandlerFunc {
	return registry[name]
}
