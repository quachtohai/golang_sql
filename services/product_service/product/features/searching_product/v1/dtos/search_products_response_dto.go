package dtos

import (
	"golang_sql/pkg/utils"
	"golang_sql/services/product_service/product/dtos"
)

type SearchProductsResponseDto struct {
	Products *utils.ListResult[*dtos.ProductDto]
}