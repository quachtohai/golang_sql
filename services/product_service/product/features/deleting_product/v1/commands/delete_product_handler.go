package commands

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/product/data/contracts"
	eventsv1 "golang_sql/services/product_service/product/features/deleting_product/v1/events"

	"github.com/mehdihadeli/go-mediatr"
)

type DeleteProductHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewDeleteProductHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *DeleteProductHandler {
	return &DeleteProductHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *DeleteProductHandler) Handle(ctx context.Context, command *DeleteProduct) (*mediatr.Unit, error) {

	if err := c.productRepository.DeleteProductByID(ctx, command.ProductID); err != nil {
		return nil, err
	}

	err := c.rabbitmqPublisher.PublishMessage(eventsv1.ProductDeleted{
		ProductId: command.ProductID,
	})
	if err != nil {
		return nil, err
	}

	c.log.Info("DeleteProduct successfully executed")

	return &mediatr.Unit{}, err
}
