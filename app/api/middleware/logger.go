package middleware

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/alfiejones/panda-ci/pkg/utils/env"
)

func Logger() echo.MiddlewareFunc {
	logger := zerolog.New(os.Stdout)

	if env.GetStage() == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogError:        true,
		LogRemoteIP:     true,
		LogMethod:       true,
		LogLatency:      true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.
				Err(v.Error).
				Str("URI", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Dur("latency", v.Latency).
				Int64("responseSize", v.ResponseSize).
				Str("remoteIP", v.RemoteIP).
				Msg("request")

			return nil
		},
	})
}
