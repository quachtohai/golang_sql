package commands

import (
	"context"
	"encoding/json"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/mapper"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/product/data/contracts"
	dtosv1 "golang_sql/services/product_service/product/features/creating_product/v1/dtos"
	eventsv1 "golang_sql/services/product_service/product/features/creating_product/v1/events"
	"golang_sql/services/product_service/product/models"
)

type CreateProductHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewCreateProductHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *CreateProductHandler {
	return &CreateProductHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *CreateProductHandler) Handle(ctx context.Context, command *CreateProduct) (*dtosv1.CreateProductResponseDto, error) {

	product := &models.Product{
		ProductId:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		InventoryId: command.InventoryId,
		Count:       command.Count,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.productRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*eventsv1.ProductCreated](createdProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(evt)
	if err != nil {
		return nil, err
	}

	response := &dtosv1.CreateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("CreateProductResponseDto", string(bytes))

	return response, nil
}
