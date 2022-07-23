package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var sensitiveAPIs = map[string]bool{}

// filterSensitiveAPI only returns `email` field for sensitive APIs
func filterSensitiveAPI(path string, data []byte) []byte {
	_, ok := sensitiveAPIs[path]
	if ok {
		return []byte{}
	}

	return data
}

// LoggerMiddleware is referenced from gin's logger implementation with additional capabilities:
// 1. use zerolog to do structure log
// 2. add requestID into context logger
func LoggerMiddleware(rootCtx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Ignore health-check to avoid polluting API logs
		if path == "/api/v1/health" {
			c.Next()
			return
		}

		// Add RequestID into the logger of the request context
		requestID := requestid.Get(c)
		zlog := zerolog.Ctx(rootCtx).With().
			Str("requestID", requestID).
			Str("path", c.FullPath()).
			Str("method", c.Request.Method).
			Logger()
		c.Request = c.Request.WithContext(zlog.WithContext(rootCtx))

		// Use TeeReader to duplicate the request body to an internal buffer, so
		// that we could use it for logging
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		c.Request.Body = io.NopCloser(tee)

		// Process request
		c.Next()

		// Build all information that we want to log
		end := time.Now()
		params := gin.LogFormatterParams{
			TimeStamp:  end,
			Latency:    end.Sub(start),
			ClientIP:   c.ClientIP(),
			StatusCode: c.Writer.Status(),
			BodySize:   c.Writer.Size(),
		}
		if err := c.Errors.Last(); err != nil {
			params.ErrorMessage = err.Error()
		}
		if raw != "" {
			path = path + "?" + raw
		}
		params.Path = path

		// Build logger with proper severity
		var l *zerolog.Event
		if params.StatusCode >= 300 || len(params.ErrorMessage) != 0 {
			l = zerolog.Ctx(c.Request.Context()).Error()
		} else {
			l = zerolog.Ctx(c.Request.Context()).Info()
		}

		l = l.Time("callTime", params.TimeStamp).
			Int("status", params.StatusCode).
			Dur("latency", params.Latency).
			Str("clientIP", params.ClientIP).
			Str("fullPath", params.Path).
			Str("component", "router").
			Str("userAgent", c.Request.Header.Get("User-Agent"))
		if params.ErrorMessage != "" {
			l = l.Err(errors.New(params.ErrorMessage))
		}
		if buf.Len() > 0 {
			data := buf.Bytes()

			// Try to filter request body if it's a sensitive API
			data = filterSensitiveAPI(params.Path, data)

			var jsonBuf bytes.Buffer
			if err := json.Compact(&jsonBuf, data); err == nil {
				l = l.RawJSON("request", jsonBuf.Bytes())
			}
		}
		l.Send()
	}
}
