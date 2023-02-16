package mwtracing

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
)

func init() {
	const semLogContext = "tracing-middleware::init"
	log.Info().Msg(semLogContext)
	mwregistry.RegisterHandlerFactory(TracingHandlerId, NewTracingHandler)
}

type TracingHandler struct {
	config *TracingHandlerConfig
}

func MustNewTracingHandler(cfg interface{}) mws.MiddlewareHandler {

	const semLogContext = "must-new-tracing-handler"
	h, err := NewTracingHandler(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	return h
}

// NewTracingHandler builds an Handler
func NewTracingHandler(cfg interface{}) (mws.MiddlewareHandler, error) {

	const semLogContext = "new-tracing-handler"
	tcfg := DefaultTracingHandlerConfig

	if cfg != nil && !reflect.ValueOf(cfg).IsNil() {
		if mapCfg, ok := cfg.(mws.MiddlewareHandlerConfig); ok {
			err := mapstructure.Decode(mapCfg, &tcfg)
			if err != nil {
				return nil, err
			}
		} else {
			log.Warn().Msg(semLogContext + " unmarshal issue for tracing handler config")
		}
	} else {
		log.Info().Str("mw-id", TracingHandlerId).Msg(semLogContext + " config null...reverting to default values")
	}

	log.Info().Str("mw-id", TracingHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	return &TracingHandler{config: &tcfg}, nil
}

func (t *TracingHandler) GetKind() string {
	return TracingHandlerKind
}

func (t *TracingHandler) HandleFunc() gin.HandlerFunc {

	return func(c *gin.Context) {

		log.Trace().Str("requestPath", c.Request.RequestURI).Send()

		var span opentracing.Span
		parentSpanCtx, serr := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if nil != serr {
			span = opentracing.StartSpan(c.FullPath())
		} else {
			span = opentracing.StartSpan(c.FullPath(), opentracing.ChildOf(parentSpanCtx))
		}
		defer span.Finish()

		t.addConfiguredTags2Span(span, c.Request)

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		if nil != c {
			c.Next()
		}

		statusCode := c.Writer.Status()
		if span != nil {
			span.SetTag("http.method", c.Request.Method)
			span.SetTag("http.status_code", statusCode)
		}

	}
}

func (t *TracingHandler) addConfiguredTags2Span(span opentracing.Span, req *http.Request) {
	for _, t := range t.config.Tags {
		switch t.Source {
		case TracingHandlerSourceTypeHeader:
			h := req.Header.Get(t.Value)
			if h != "" {
				span.SetTag(t.Name, h)
			}
		}
	}
}

//func (t *TracingHandler) fail(c *gin.Context, appErr AppError, span opentracing.Span) {
//
//	if nil != span {
//		ext.Error.Set(span, true)
//		span.SetTag("cause", appErr)
//		ext.HTTPStatusCode.Set(span, uint16(appErr.GetCode()))
//	}
//
//	// injecting error id and tagging span
//	errid, err := gonanoid.Generate(t.config.Alphabet, 32)
//	if nil != err { // in this case just dump error, we want error handling to be smooth
//		// ignore
//	} else {
//		if nil != span {
//			span.SetTag(t.config.SpanTag, errid)
//			c.Header(t.config.Header, errid)
//		}
//	}
//}
//
//func (t *TracingHandler) failWithContext(c *gin.Context, w http.ResponseWriter, appErr AppError) {
//	span := opentracing.SpanFromContext(c.Request.Context())
//	t.fail(c, appErr, span)
//}
