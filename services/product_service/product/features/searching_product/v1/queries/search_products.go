package queries

import (
	"golang_sql/pkg/utils"
)

type SearchProducts struct {
	SearchText string `validate:"required"`
	*utils.ListQuery
}
