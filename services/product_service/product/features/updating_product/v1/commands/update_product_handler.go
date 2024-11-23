package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/mapper"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/product/data/contracts"
	dtosv1 "golang_sql/services/product_service/product/features/updating_product/v1/dtos"
	eventsv1 "golang_sql/services/product_service/product/features/updating_product/v1/events"
	"golang_sql/services/product_service/product/models"

	"github.com/pkg/errors"
)

type UpdateProductHandler struct {
	log               logger.ILogger
	rabbitmqPublisher rabbitmq.IPublisher
	productRepository contracts.ProductRepository
	ctx               context.Context
}

func NewUpdateProductHandler(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	productRepository contracts.ProductRepository, ctx context.Context) *UpdateProductHandler {
	return &UpdateProductHandler{log: log, productRepository: productRepository, ctx: ctx, rabbitmqPublisher: rabbitmqPublisher}
}

func (c *UpdateProductHandler) Handle(ctx context.Context, command *UpdateProduct) (*dtosv1.UpdateProductResponseDto, error) {

	//simple call grpcClient
	//identityGrpcClient := identity_service.NewIdentityServiceClient(c.grpcClient.GetGrpcConnection())
	//user, err := identityGrpcClient.GetUserById(ctx, &identity_service.GetUserByIdReq{UserId: "1"})
	//if err != nil {
	//	return nil, err
	//}
	//
	//c.log.Infof("userId: %s", user.User.UserId)

	_, err := c.productRepository.GetProductById(ctx, command.ProductID)

	if err != nil {
		notFoundErr := errors.Wrap(err, fmt.Sprintf("product with id %s not found", command.ProductID))
		return nil, notFoundErr
	}

	product := &models.Product{ProductId: command.ProductID, Name: command.Name, Description: command.Description, Price: command.Price, UpdatedAt: command.UpdatedAt, InventoryId: command.InventoryId, Count: command.Count}

	updatedProduct, err := c.productRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	evt, err := mapper.Map[*eventsv1.ProductUpdated](updatedProduct)
	if err != nil {
		return nil, err
	}

	err = c.rabbitmqPublisher.PublishMessage(evt)

	response := &dtosv1.UpdateProductResponseDto{ProductId: product.ProductId}
	bytes, _ := json.Marshal(response)

	c.log.Info("UpdateProductResponseDto", string(bytes))

	return response, nil
}
