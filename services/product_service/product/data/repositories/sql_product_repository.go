package repositories

import (
	"context"
	"fmt"
	gormsql "golang_sql/pkg/gorm_sql"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/utils"
	"golang_sql/services/product_service/product/data/contracts"
	"golang_sql/services/product_service/product/models"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type SqlProductRepository struct {
	log  logger.ILogger
	cfg  *gormsql.GormSqlConfig
	gorm *gorm.DB
}

func NewSqlProductRepository(log logger.ILogger, cfg *gormsql.GormSqlConfig, gorm *gorm.DB) contracts.ProductRepository {
	return &SqlProductRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *SqlProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	result, err := gormsql.Paginate[*models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *SqlProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gormsql.Paginate[*models.Product](ctx, listQuery, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *SqlProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	return &product, nil
}

func (p *SqlProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return product, nil
}

func (p *SqlProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error in updating product with id %s into the database.", updateProduct.ProductId))
	}

	return updateProduct, nil
}

func (p *SqlProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return errors.Wrap(err, "error in the deleting product into the database.")
	}

	return nil
}
