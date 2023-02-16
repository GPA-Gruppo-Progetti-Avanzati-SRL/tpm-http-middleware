package mws

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareHandlerConfig map[string]interface{}

/*
struct {
	ErrCfg     *ErrorHandlerConfig           `yaml:"gin-mw-error" mapstructure:"gin-mw-error" json:"gin-mw-error"`
	MetricsCfg *PromHttpMetricsHandlerConfig `yaml:"gin-mw-metrics" mapstructure:"gin-mw-metrics" json:"gin-mw-metrics"`
	TraceCfg   *TracingHandlerConfig         `yaml:"gin-mw-tracing" mapstructure:"gin-mw-tracing" json:"gin-mw-tracing"`
}
*/

type MiddlewareHandler interface {
	GetKind() string
	HandleFunc() gin.HandlerFunc
}

/*
 * Package Configuration defaults

func GetConfigDefaults(contextPath string) []configuration.VarDefinition {
	return []configuration.VarDefinition{
		{strings.Join([]string{contextPath, ErrorHandlerId, "with-cause"}, "."), ErrorHandlerDefaultWithCause, "error is in clear"},
		{strings.Join([]string{contextPath, ErrorHandlerId, "alphabet"}, "."), ErrorHandlerDefaultAlphabet, "alphabet"},
		{strings.Join([]string{contextPath, ErrorHandlerId, "spantag"}, "."), ErrorHandlerDefaultSpanTag, "spantag"},
		{strings.Join([]string{contextPath, ErrorHandlerId, "header"}, "."), TErrorHandlerDefaultHeader, "header"},
	}
}
*/
