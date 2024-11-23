package repositories

import (
	"context"
	"fmt"
	gormsql "golang_sql/pkg/gorm_sql"
	"golang_sql/pkg/logger"
	contracts "golang_sql/services/inventory_service/inventory/data/contracts"
	"golang_sql/services/inventory_service/inventory/models"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type SqlInventoryRepository struct {
	log  logger.ILogger
	cfg  *gormsql.GormSqlConfig
	gorm *gorm.DB
}

func NewSqlInventoryRepository(log logger.ILogger, cfg *gormsql.GormSqlConfig, gorm *gorm.DB) contracts.InventoryRepository {
	return &SqlInventoryRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *SqlInventoryRepository) AddProductItemToInventory(ctx context.Context, productItem *models.ProductItem) (*models.ProductItem, error) {

	if err := p.gorm.Create(&productItem).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return productItem, nil
}

func (p *SqlInventoryRepository) GetProductInInventories(ctx context.Context, uuid uuid.UUID) (*models.ProductItem, error) {
	var productItem models.ProductItem

	if err := p.gorm.First(&productItem, uuid).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product item with id %s into the database.", uuid))
	}

	return &productItem, nil
}
