package queries

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/pkg/utils"
	"golang_sql/services/product_service/product/data/contracts"
	"golang_sql/services/product_service/product/dtos"
	dtosv1 "golang_sql/services/product_service/product/features/searching_product/v1/dtos"
)

type SearchProductsHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewSearchProductsHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *SearchProductsHandler {
	return &SearchProductsHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *SearchProductsHandler) Handle(ctx context.Context, query *SearchProducts) (*dtosv1.SearchProductsResponseDto, error) {

	products, err := c.productRepository.SearchProducts(ctx, query.SearchText, query.ListQuery)
	if err != nil {
		return nil, err
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtos.ProductDto](products)
	if err != nil {
		return nil, err
	}

	return &dtosv1.SearchProductsResponseDto{Products: listResultDto}, nil
}
