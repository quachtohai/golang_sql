package endpoints

import (
	"context"
	echomiddleware "golang_sql/pkg/http/echo/middleware"
	"golang_sql/pkg/logger"
	commandsv1 "golang_sql/services/product_service/product/features/updating_product/v1/commands"
	dtosv1 "golang_sql/services/product_service/product/features/updating_product/v1/dtos"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.PUT("/:id", updateProduct(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

// UpdateProduct
// @Tags        Products
// @Summary     Update product
// @Description Update existing product
// @Accept      json
// @Produce     json
// @Param       UpdateProductRequestDto body dtos.UpdateProductRequestDto true "Product data"
// @Param       id                      path string                       true "Product ID"
// @Success     204
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [put]
func updateProduct(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := &dtosv1.UpdateProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[updateProductEndpoint_handler.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commandsv1.NewUpdateProduct(request.ProductId, request.Name, request.Description, request.Price, request.InventoryId, request.Count)

		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[updateProductEndpoint_handler.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		_, err := mediatr.Send[*commandsv1.UpdateProduct, *dtosv1.UpdateProductResponseDto](ctx, command)

		if err != nil {
			log.Warnf("UpdateProduct", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("(product updated) id: {%s}", request.ProductId)

		return c.NoContent(http.StatusNoContent)
	}
}
