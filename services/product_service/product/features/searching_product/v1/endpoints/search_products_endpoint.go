package endpoints

import (
	"context"
	echomiddleware "golang_sql/pkg/http/echo/middleware"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/utils"
	dtosv1 "golang_sql/services/product_service/product/features/searching_product/v1/dtos"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/products")
	group.GET("/search", searchProducts(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

// SearchProducts
// @Tags        Products
// @Summary     Search products
// @Description Search products
// @Accept      json
// @Produce     json
// @Param       searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success     200  {object} dtos.SearchProductsResponseDto
// @Security ApiKeyAuth
// @Router      /api/v1/products/search [get]
func searchProducts(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {

		listQuery, err := utils.GetListQueryFromCtx(c)

		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		request := &dtosv1.SearchProductsRequestDto{ListQuery: listQuery}

		// https://echo.labstack.com/guide/binding/
		if err := c.Bind(request); err != nil {
			log.Warn("Bind", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		query := &dtosv1.SearchProductsRequestDto{SearchText: request.SearchText, ListQuery: request.ListQuery}

		if err := validator.StructCtx(ctx, query); err != nil {
			log.Errorf("(validate) err: {%v}", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		queryResult, err := mediatr.Send[*dtosv1.SearchProductsRequestDto, *dtosv1.SearchProductsResponseDto](ctx, query)

		if err != nil {
			log.Warn("SearchProducts", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
