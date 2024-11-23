package queries

import (
	"context"
	"fmt"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/mapper"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/product/data/contracts"
	"golang_sql/services/product_service/product/dtos"
	dtosv1 "golang_sql/services/product_service/product/features/getting_product_by_id/v1/dtos"

	"github.com/pkg/errors"
)

type GetProductByIdHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewGetProductByIdHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *GetProductByIdHandler {
	return &GetProductByIdHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (q *GetProductByIdHandler) Handle(ctx context.Context, query *GetProductById) (*dtosv1.GetProductByIdResponseDto, error) {

	product, err := q.productRepository.GetProductById(ctx, query.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", query.ProductID))
		return nil, notFoundErr
	}

	productDto, err := mapper.Map[*dtos.ProductDto](product)
	if err != nil {
		return nil, err
	}

	return &dtosv1.GetProductByIdResponseDto{Product: productDto}, nil
}
