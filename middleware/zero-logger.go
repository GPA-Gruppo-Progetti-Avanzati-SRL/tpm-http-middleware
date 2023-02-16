package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

/*
 * Version adapted from https://github.com/dn365/gin-zerolog/blob/master/gin_zerolog.go
 */

type ginHands struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
}

func ZeroLogger(serName string, notlogged ...string) gin.HandlerFunc {

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()

		if _, ok := skip[path]; !ok {
			// after request
			// latency := time.Since(t)
			// clientIP := c.ClientIP()
			// method := c.Request.Method
			// statusCode := c.Writer.Status()
			if raw != "" {
				path = path + "?" + raw
			}
			msg := c.Errors.String()
			if msg == "" {
				msg = "Request"
			}
			cData := &ginHands{
				SerName:    serName,
				Path:       path,
				Latency:    time.Since(t),
				Method:     c.Request.Method,
				StatusCode: c.Writer.Status(),
				ClientIP:   c.ClientIP(),
				MsgStr:     msg,
			}

			logSwitch(cData)
		}
	}
}

func logSwitch(data *ginHands) {
	switch {
	case data.StatusCode >= 400 && data.StatusCode < 500:
		{
			log.Warn().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	case data.StatusCode >= 500:
		{
			log.Error().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
		}
	default:
		log.Info().Str("ser_name", data.SerName).Str("method", data.Method).Str("path", data.Path).Dur("resp_time", data.Latency).Int("status", data.StatusCode).Str("client_ip", data.ClientIP).Msg(data.MsgStr)
	}
}
