package configurations

import (
	"context"
	"golang_sql/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ConfigEndpoints(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {

}
