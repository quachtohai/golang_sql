package configurations

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/services/identity_service/identity/features/registering_user/v1/endpoints"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ConfigEndpoints(validator *validator.Validate, logger logger.ILogger, echo *echo.Echo, ctx context.Context) {
	endpoints.MapRoute(validator, logger, echo, ctx)
}
