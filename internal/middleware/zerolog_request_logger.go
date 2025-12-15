package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/th3oth3rjak3/mainframe/internal/shared"
)

func ZerologRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			var finalStatus int
			if err != nil {
				finalStatus, _ = shared.ResolveError(err)
			} else {
				finalStatus = res.Status
			}

			duration := stop.Sub(start)
			latencyMs := float64(duration.Nanoseconds()) / 1_000_000.0

			event := log.Info() // Default to Info
			if finalStatus >= 500 {
				event = log.Error()
			} else if finalStatus >= 400 {
				event = log.Warn()
			}

			// Log using Zerolog
			event.
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", finalStatus).
				Str("remote_ip", c.RealIP()).
				Float64("latency_ms", latencyMs).
				Str("user_agent", req.UserAgent()).
				Msg("request")

			return err
		}
	}
}
