package mwerror

const (
	ErrorHandlerId   = "gin-mw-error"
	ErrorHandlerKind = "mw-kind-error"

	ErrorHandlerDefaultAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-"
	ErrorHandlerDefaultSpanTag  = "error.id"
	ErrorHandlerDefaultHeader   = "x-errid"

	ErrorHandlerCauseModeUseFirst = "use-first"
	ErrorHandlerCauseModeUseLast  = "use-last"
	ErrorHandlerCauseModeNone     = "none"
	ErrorHandlerDefaultWithCause  = true

	ErrorHandlerStatusCodeHandlingPolicyModeNone       = "none"
	ErrorHandlerStatusCodeHandlingPolicyModeIfListed   = "if-listed"
	ErrorHandlerStatusCodeHandlingPolicyModeIfUnListed = "if-unlisted"
	ErrorHandlerDefaultStatusCodeHandlingPolicyMode    = ErrorHandlerStatusCodeHandlingPolicyModeNone
)

type StatusCodeRange struct {
	StatusCode     int `yaml:"value,omitempty"  mapstructure:"value,omitempty"  json:"value,omitempty"`
	StatusCodeFrom int `yaml:"from,omitempty"  mapstructure:"from,omitempty"  json:"from,omitempty"`
	StatusCodeTo   int `yaml:"to,omitempty"  mapstructure:"to,omitempty"  json:"to,omitempty"`
}

func (r StatusCodeRange) In(c int) bool {
	if c == 0 {
		return false
	}

	if c == r.StatusCode {
		return true
	}

	if r.StatusCodeFrom > 0 || r.StatusCodeTo > 0 {
		if r.StatusCodeFrom > 0 && c < r.StatusCodeFrom {
			return false
		}

		if r.StatusCodeTo > 0 && c > r.StatusCodeTo {
			return false
		}

		return true
	}

	return false
}

type StatusCodeHandlingPolicy struct {
	PolicyMode       string            `yaml:"policy-mode,omitempty"  mapstructure:"policy-mode,omitempty"  json:"policy-mode,omitempty"`
	StatusCodeRanges []StatusCodeRange `yaml:"status-codes,omitempty"  mapstructure:"status-codes,omitempty"  json:"status-codes,omitempty"`
}

func (p StatusCodeHandlingPolicy) Hightlight(c int) bool {

	if p.PolicyMode == "" || p.PolicyMode == ErrorHandlerStatusCodeHandlingPolicyModeNone {
		return false
	}

	incl := false
	for _, r := range p.StatusCodeRanges {
		if r.In(c) {
			incl = true
			break
		}
	}

	hl := false
	switch p.PolicyMode {
	case ErrorHandlerStatusCodeHandlingPolicyModeIfListed:
		hl = incl
	case ErrorHandlerStatusCodeHandlingPolicyModeIfUnListed:
		hl = !incl
	}

	return hl
}

//    WithErrorEnabled(bool)   Enables/Disables error disclosure to the client
//                             if enabled the http error description is propagated to the client
//                             if disabled a response Header, configured with WithErrorDisclosureHeader is returned
//                             to the client with an errorid and the error is injected in an opentracing span having
//                             the same id as tag
//    WithSpanTag(string)      span tag for the error  (defaults to "error.id")
//    WithHeader(string)       error id header (defaults to "x-errid")
//    WithAlphabet(string)     alphabet  to generate the error id

type ErrorHandlerConfig struct {
	WithCause                bool                     `yaml:"with-cause,omitempty"  mapstructure:"with-cause,omitempty" json:"with-cause,omitempty"`
	Alphabet                 string                   `yaml:"alphabet,omitempty"  mapstructure:"alphabet,omitempty"  json:"alphabet,omitempty"`
	SpanTag                  string                   `yaml:"span-tag,omitempty"  mapstructure:"span-tag,omitempty"  json:"span-tag,omitempty"`
	Header                   string                   `yaml:"header,omitempty"  mapstructure:"header,omitempty"  json:"header,omitempty"`
	StatusCodeHandlingPolicy StatusCodeHandlingPolicy `yaml:"status-code-policy,omitempty"  mapstructure:"status-code-policy,omitempty"  json:"status-code-policy,omitempty"`
}

var DefaultErrorHandlerConfig = ErrorHandlerConfig{
	WithCause: ErrorHandlerDefaultWithCause,
	Alphabet:  ErrorHandlerDefaultAlphabet,
	SpanTag:   ErrorHandlerDefaultSpanTag,
	Header:    ErrorHandlerDefaultHeader,
	StatusCodeHandlingPolicy: StatusCodeHandlingPolicy{
		PolicyMode: ErrorHandlerDefaultStatusCodeHandlingPolicyMode,
	},
}

func (h *ErrorHandlerConfig) GetKind() string {
	return ErrorHandlerKind
}

type ErrorHandlerConfigOption func(*ErrorHandlerConfig)
type ErrorHandlerConfigBuilder struct {
	opts []ErrorHandlerConfigOption
}

func (cb *ErrorHandlerConfigBuilder) WithCauseMode(enabled bool) *ErrorHandlerConfigBuilder {

	f := func(c *ErrorHandlerConfig) {
		c.WithCause = enabled
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *ErrorHandlerConfigBuilder) WithAlphabet(alphabet string) *ErrorHandlerConfigBuilder {
	f := func(c *ErrorHandlerConfig) {
		c.Alphabet = alphabet
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *ErrorHandlerConfigBuilder) WithSpanTag(s string) *ErrorHandlerConfigBuilder {
	f := func(c *ErrorHandlerConfig) {
		c.SpanTag = s
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *ErrorHandlerConfigBuilder) WithHeader(h string) *ErrorHandlerConfigBuilder {
	f := func(c *ErrorHandlerConfig) {
		c.Header = h
	}

	cb.opts = append(cb.opts, f)
	return cb
}

func (cb *ErrorHandlerConfigBuilder) Build() *ErrorHandlerConfig {
	c := DefaultErrorHandlerConfig

	for _, o := range cb.opts {
		o(&c)
	}

	return &c
}
