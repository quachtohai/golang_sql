package contracts

import (
	"context"
	"golang_sql/services/inventory_service/inventory/models"

	uuid "github.com/satori/go.uuid"
)

type InventoryRepository interface {
	AddProductItemToInventory(ctx context.Context, inventory *models.ProductItem) (*models.ProductItem, error)
	GetProductInInventories(ctx context.Context, productId uuid.UUID) (*models.ProductItem, error)
}
