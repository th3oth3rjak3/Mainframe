package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func ZerologRequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()
			duration := stop.Sub(start)
			latencyMs := float64(duration.Nanoseconds()) / 1_000_000.0
			// Log using Zerolog
			log.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Str("remote_ip", c.RealIP()).
				Float64("latency_ms", latencyMs).
				Str("user_agent", req.UserAgent()).
				Msg("request")

			return err
		}
	}
}
