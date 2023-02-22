package mwregistry

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HandlerCatalogConfig map[string]interface{}
type HandlerFactory func(interface{}) (mws.MiddlewareHandler, error)
type HandlerRegistry map[string]gin.HandlerFunc

var handlerFactoryMap map[string]HandlerFactory

func RegisterHandlerFactory(handlerId string, hf HandlerFactory) {
	const semLogContext = "middleware:register-handler-factory"

	if handlerFactoryMap == nil {
		handlerFactoryMap = make(map[string]HandlerFactory)
	}

	if _, ok := handlerFactoryMap[handlerId]; ok {
		log.Warn().Str("mw-id", handlerId).Msg(semLogContext + " handler factory already registered")
		return
	}

	handlerFactoryMap[handlerId] = hf
}

var registry HandlerRegistry = make(map[string]gin.HandlerFunc)

func InitializeHandlerRegistry(registryConfig HandlerCatalogConfig, mwInUse []string) error {

	const semLogContext = "middleware::registry-initialization"

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

func GetMiddlewareHandler(name string) gin.HandlerFunc {
	return registry[name]
}
