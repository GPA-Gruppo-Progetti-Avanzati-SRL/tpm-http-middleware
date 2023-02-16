package middleware

import (
	"github.com/gin-gonic/gin"
)

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
