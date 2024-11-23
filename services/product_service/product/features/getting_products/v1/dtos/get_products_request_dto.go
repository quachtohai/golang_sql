package dtos

import (
	"golang_sql/pkg/utils"
)

type GetProductsRequestDto struct {
	*utils.ListQuery
}
