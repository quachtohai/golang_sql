package configurations

import (
	"context"
	"golang_sql/pkg/logger"
	creating_product "golang_sql/services/product_service/product/features/creating_product/v1/endpoints"
	deleting_product "golang_sql/services/product_service/product/features/deleting_product/v1/endpoints"
	getting_product_by_id "golang_sql/services/product_service/product/features/getting_product_by_id/v1/endpoints"
	getting_products "golang_sql/services/product_service/product/features/getting_products/v1/endpoints"
	searching_product "golang_sql/services/product_service/product/features/searching_product/v1/endpoints"
	updating_product "golang_sql/services/product_service/product/features/updating_product/v1/endpoints"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func ConfigEndpoints(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {

	creating_product.MapRoute(validator, log, echo, ctx)
	deleting_product.MapRoute(validator, log, echo, ctx)
	getting_product_by_id.MapRoute(validator, log, echo, ctx)
	getting_products.MapRoute(validator, log, echo, ctx)
	searching_product.MapRoute(validator, log, echo, ctx)
	updating_product.MapRoute(validator, log, echo, ctx)
}
