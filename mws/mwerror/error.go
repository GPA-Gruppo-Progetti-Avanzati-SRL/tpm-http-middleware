package mwerror

import (
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/rs/zerolog/log"
	"net/http"
	"reflect"
)

func init() {
	const semLogContext = "error-middleware::init"
	log.Info().Msg(semLogContext)
	mwregistry.RegisterHandlerFactory(ErrorHandlerId, NewErrorHandler)
}

type ErrorHandler struct {
	config *ErrorHandlerConfig
}

func MustNewErrorHandler(cfg interface{}) mws.MiddlewareHandler {

	const semLogContext = "must-new-error-handler"
	h, err := NewErrorHandler(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	return h
}

func NewErrorHandler(cfg interface{}) (mws.MiddlewareHandler, error) {

	const semLogContext = "error-handler::new"

	tcfg := DefaultErrorHandlerConfig
	if cfg != nil && !reflect.ValueOf(cfg).IsNil() {
		switch typedCfg := cfg.(type) {
		case mws.MiddlewareHandlerConfig:
			err := mapstructure.Decode(typedCfg, &tcfg)
			if err != nil {
				return nil, err
			}

			log.Info().Interface("cfg", &tcfg).Msg(semLogContext)
		default:
			log.Warn().Msg(semLogContext + " unmarshal issue for error handler config")
		}
	} else {
		log.Info().Str("mw-id", ErrorHandlerId).Msg(semLogContext + " config null...reverting to default values")
	}

	log.Info().Str("mw-id", ErrorHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	return &ErrorHandler{config: &tcfg}, nil
}

func (h *ErrorHandler) GetKind() string {
	return ErrorHandlerKind
}

func (h *ErrorHandler) HandleFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if nil != c {
			c.Next()

			if len(c.Errors) > 0 {
				for _, e := range c.Errors {
					log.Error().Str("middleware", "error").Msg(e.Error())
				}

				var ae AppError
				ae = getAppError(c.Errors[0])
				if !h.config.WithCause {
					ae = ae.Sanitized()
				}

				span := opentracing.SpanFromContext(c.Request.Context())
				if nil != span {
					h.fail(c, span, ae.Error())
				}

				c.AbortWithStatusJSON(ae.GetStatusCode(), ae)
			} else if h.config.StatusCodeHandlingPolicy.Highlight(c.Writer.Status()) {
				span := opentracing.SpanFromContext(c.Request.Context())
				if nil != span {
					// In this case headers have been written already.... c.Writer.Written() == true
					// limit to setting trace span error flag
					ext.Error.Set(span, true)
				}
			}
		}
	}
}

func getAppError(err error) AppError {
	var parsedError AppError

	gerr, ok1 := err.(*gin.Error)
	if !ok1 {
		parsedError = &AppErrorImpl{
			StatusCode: http.StatusInternalServerError,
			Text:       "Internal Server Error",
		}

		return parsedError
	}

	switch v := gerr.Err.(type) {
	case AppError:
		parsedError = v
	default:
		parsedError = &AppErrorImpl{
			StatusCode: http.StatusInternalServerError,
			Text:       "Internal Server Error",
			Message:    v.Error(),
		}
	}

	return parsedError
}

func (h *ErrorHandler) fail(c *gin.Context, span opentracing.Span, cause string) {
	ext.Error.Set(span, true)
	if cause != "" {
		span.SetTag("cause", cause)
	}

	// injecting error id and tagging span
	errid, err := gonanoid.Generate(h.config.Alphabet, 32)
	if nil != err {
		// in this case just dump error, we want error handling to be smooth
		// ignore
	} else {
		if nil != span {
			span.SetTag(h.config.SpanTag, errid)
			c.Header(h.config.Header, errid)
		}
	}
}
