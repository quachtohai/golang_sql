package dtos

import (
	"golang_sql/pkg/utils"
)

type SearchProductsRequestDto struct {
	SearchText       string `query:"search" json:"search"`
	*utils.ListQuery `json:"listQuery"`
}
