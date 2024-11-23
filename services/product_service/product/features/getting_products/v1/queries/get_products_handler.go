package queries

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/pkg/utils"
	"golang_sql/services/product_service/product/data/contracts"
	"golang_sql/services/product_service/product/dtos"
	dtosv1 "golang_sql/services/product_service/product/features/getting_products/v1/dtos"
)

type GetProductsHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewGetProductsHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *GetProductsHandler {
	return &GetProductsHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *GetProductsHandler) Handle(ctx context.Context, query *GetProducts) (*dtosv1.GetProductsResponseDto, error) {

	products, err := c.productRepository.GetAllProducts(ctx, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)

	if err != nil {
		return nil, err
	}
	return &dtosv1.GetProductsResponseDto{Products: listResultDto}, nil
}
