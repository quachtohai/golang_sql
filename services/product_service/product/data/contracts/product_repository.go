package contracts

import (
	"context"
	"golang_sql/pkg/utils"
	"golang_sql/services/product_service/product/models"

	uuid "github.com/satori/go.uuid"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error)
	SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error)
	GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProductByID(ctx context.Context, uuid uuid.UUID) error
}
