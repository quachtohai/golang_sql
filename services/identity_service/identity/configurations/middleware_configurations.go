package configurations

import (
	"strings"

	echomiddleware "golang_sql/pkg/http/echo/middleware"
	"golang_sql/pkg/otel"
	otelmiddleware "golang_sql/pkg/otel/middleware"
	"golang_sql/services/identity_service/identity/constants"
	"golang_sql/services/identity_service/identity/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigMiddlewares(e *echo.Echo, jaegerCfg *otel.JaegerConfig) {

	e.HideBanner = false

	e.Use(middleware.Logger())
	e.HTTPErrorHandler = middlewares.ProblemDetailsHandler
	e.Use(otelmiddleware.EchoTracerMiddleware(jaegerCfg.ServiceName))

	e.Use(echomiddleware.CorrelationIdMiddleware)
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: constants.GzipLevel,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))

	e.Use(middleware.BodyLimit(constants.BodyLimit))
}
