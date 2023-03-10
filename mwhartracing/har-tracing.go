package mwhartracing

import (
	"bytes"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/har"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-archive/hartracing"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mwregistry"
	"github.com/GPA-Gruppo-Progetti-Avanzati-SRL/tpm-http-middleware/mws"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

func init() {
	const semLogContext = "har-tracing-middleware::init"
	log.Info().Msg(semLogContext)
	mwregistry.RegisterHandlerFactory(HarTracingHandlerId, NewHarTracingHandler)
}

type HarTracingHandler struct {
	config *HarTracingHandlerConfig
}

type bodyBufferedWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyBufferedWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func MustNewTracingHandler(cfg interface{}) mws.MiddlewareHandler {

	const semLogContext = "har-tracing-handler::must-new"
	h, err := NewHarTracingHandler(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg(semLogContext)
	}

	return h
}

// NewHarTracingHandler builds an Handler
func NewHarTracingHandler(cfg interface{}) (mws.MiddlewareHandler, error) {

	const semLogContext = "har-tracing-handler::new"
	tcfg := DefaultTracingHandlerConfig

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
		case *HarTracingHandlerConfig:
		default:
			log.Warn().Msg(semLogContext + " unmarshal issue for tracing handler config")
		}

	} else {
		log.Info().Str("mw-id", HarTracingHandlerId).Msg(semLogContext + " config null...reverting to default values")
	}

	log.Info().Str("mw-id", HarTracingHandlerId).Interface("cfg", tcfg).Msg(semLogContext + " handler loaded config")

	return &HarTracingHandler{config: &tcfg}, nil
}

func (t *HarTracingHandler) GetKind() string {
	return HarTracingHandlerKind
}

func (t *HarTracingHandler) HandleFunc() gin.HandlerFunc {

	const semLogContext = "har-tracing-middleware::handle"

	return func(c *gin.Context) {
		log.Trace().Str("requestPath", c.Request.RequestURI).Msg(semLogContext)

		var harSpan hartracing.Span
		var entry har.Entry
		parentSpanCtx, serr := hartracing.GlobalTracer().Extract("", hartracing.HTTPHeadersCarrier(c.Request.Header))
		if nil != serr {
			// No incoming harSpan. Need to create a new root one with an actual entry.
			harSpan = hartracing.GlobalTracer().StartSpan()
			log.Trace().Str("harSpan-id", harSpan.Id()).Msg(semLogContext + " - starting a brand new harSpan")
			blw := &bodyBufferedWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
			entry = getRequestEntry(c)

			span := opentracing.SpanFromContext(c.Request.Context())
			if nil != span {
				span.SetTag(hartracing.HARTraceOpenTracingTagName, harSpan.Id())
			}

		} else {
			harSpan = hartracing.GlobalTracer().StartSpan(hartracing.ChildOf(parentSpanCtx))
			log.Trace().Str("harSpan-id", harSpan.Id()).Str("parent-harSpan-id", parentSpanCtx.Id()).Msg(semLogContext + " - started a child harSpan")
		}
		defer harSpan.Finish()

		c.Request = c.Request.WithContext(hartracing.ContextWithSpan(c.Request.Context(), harSpan))

		if nil != c {
			c.Next()
		}

		if entry.Request != nil {
			getResponseEntry(c, &entry)
			harSpan.AddEntry(&entry)
		}

		log.Trace().Msg(semLogContext)
	}
}

func getRequestEntry(c *gin.Context) har.Entry {
	var hs har.NameValuePairs
	var ct string
	for n, h := range c.Request.Header {
		if strings.ToLower(n) == "content-type" {
			ct = getFirstHeaderValue(h)
		}
		for _, v := range h {
			hs = append(hs, har.NameValuePair{Name: n, Value: v})
		}
	}

	var postData *har.PostData
	bodySize := -1
	if c.Request.ContentLength > 0 {
		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		if len(body) > 0 {
			bodySize = len(body)
			postData = &har.PostData{
				MimeType: ct,
				Data:     body,
				Params:   []har.Param{},
			}
		}
	}

	req := &har.Request{
		Method:      c.Request.Method,
		URL:         c.Request.RequestURI,
		HTTPVersion: "1.1",
		Headers:     hs,
		HeadersSize: -1,
		Cookies:     []har.Cookie{},
		QueryString: []har.NameValuePair{},
		BodySize:    int64(bodySize),
		PostData:    postData,
	}

	now := time.Now()
	e := har.Entry{
		Comment:         "",
		StartedDateTime: now.Format("2006-01-02T15:04:05.999999999Z07:00"),
		StartDateTimeTm: now,
		Request:         req,
	}

	return e
}

func getResponseEntry(c *gin.Context, e *har.Entry) {

	bw, ok := c.Writer.(*bodyBufferedWriter)
	if !ok {
		return
	}

	r := &har.Response{
		Status:      c.Writer.Status(),
		HTTPVersion: "1.1",
		StatusText:  http.StatusText(c.Writer.Status()),
		HeadersSize: -1,
		BodySize:    int64(c.Writer.Size()),
		Cookies:     []har.Cookie{},
		Content: &har.Content{
			MimeType: c.Writer.Header().Get("Content-type"),
			Size:     int64(c.Writer.Size()),
			Data:     bw.body.Bytes(),
		},
	}

	for n, _ := range c.Writer.Header() {
		r.Headers = append(r.Headers, har.NameValuePair{Name: n, Value: c.Writer.Header().Get(n)})
	}

	if e.StartedDateTime != "" {
		elapsed := time.Since(e.StartDateTimeTm)
		e.Time = float64(elapsed.Milliseconds())
	}

	e.Timings = &har.Timings{
		Blocked: -1,
		DNS:     -1,
		Connect: -1,
		Send:    -1,
		Wait:    e.Time,
		Receive: -1,
		Ssl:     -1,
	}

	e.Response = r
}

func getFirstHeaderValue(vals []string) string {
	if len(vals) > 0 {
		return vals[0]
	}

	return ""
}
