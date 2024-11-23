package endpoints

import (
	"context"
	echomiddleware "golang_sql/pkg/http/echo/middleware"
	"golang_sql/pkg/logger"
	dtosv1 "golang_sql/services/product_service/product/features/getting_product_by_id/v1/dtos"
	queriesv1 "golang_sql/services/product_service/product/features/getting_product_by_id/v1/queries"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.GET("/:id", getProductByID(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

// GetProductByID
// @Tags        Products
// @Summary     Get product
// @Description Get product by id
// @Accept      json
// @Produce     json
// @Param       id  path     string true "Product ID"
// @Success     200 {object} dtos.GetProductByIdResponseDto
// @Security ApiKeyAuth
// @Router      /api/v1/products/{id} [get]
func getProductByID(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		request := &dtosv1.GetProductByIdRequestDto{}
		if err := c.Bind(request); err != nil {
			log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := queriesv1.NewGetProductById(request.ProductId)

		if err := validator.StructCtx(ctx, query); err != nil {
			log.Warn("validate", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		queryResult, err := mediatr.Send[*queriesv1.GetProductById, *dtosv1.GetProductByIdResponseDto](ctx, query)

		if err != nil {
			log.Warn("GetProductById", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
