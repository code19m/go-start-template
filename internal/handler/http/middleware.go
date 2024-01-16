package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func accessLoggerMiddleware(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the start time of the request
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		method := c.Request.Method

		// Call the next middleware in the chain
		c.Next()

		// Record the end time of the request
		end := time.Now()
		duration := end.Sub(start)
		duration = duration.Truncate(time.Microsecond)

		if duration > time.Minute {
			duration = duration.Truncate(time.Millisecond)
		}

		statusCode := c.Writer.Status()

		log := log.
			With("handler", "http").
			With("method", method).
			With("path", path).
			With("duration", duration).
			With("status", statusCode).
			With("client", c.ClientIP())

		switch {
		case statusCode >= 500:
			errMsg := ""
			ginErr := c.Errors.Last()
			if ginErr != nil {
				errMsg = ginErr.Error()
			}
			log.Error(errMsg)
		case statusCode >= 400:
			log.Warn("")
		default:
			log.Info("")
		}
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusOK)
		}
	}
}
