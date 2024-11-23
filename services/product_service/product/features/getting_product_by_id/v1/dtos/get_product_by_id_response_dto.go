package dtos

import (
	"golang_sql/services/product_service/product/dtos"
)

type GetProductByIdResponseDto struct {
	Product *dtos.ProductDto `json:"product"`
}
