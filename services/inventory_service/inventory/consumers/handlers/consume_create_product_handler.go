package handlers

import (
	"encoding/json"

	"golang_sql/pkg/mapper"
	"golang_sql/services/inventory_service/inventory/consumers/events"
	"golang_sql/services/inventory_service/inventory/models"
	"golang_sql/services/inventory_service/shared/delivery"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func HandleConsumeCreateProduct(queue string, msg amqp.Delivery, inventoryDeliveryBase *delivery.InventoryDeliveryBase) error {

	log.Infof("Message received on queue: %s with message: %s", queue, string(msg.Body))
	log.Infof("Quách Tố Hải")
	var productCreated events.ProductCreated

	err := json.Unmarshal(msg.Body, &productCreated)
	if err != nil {
		return err
	}

	count := productCreated.Count

	productItem, _ := inventoryDeliveryBase.InventoryRepository.GetProductInInventories(inventoryDeliveryBase.Ctx, productCreated.ProductId)

	if productItem != nil {
		count = productItem.Count + count
	}

	p, err := inventoryDeliveryBase.InventoryRepository.AddProductItemToInventory(inventoryDeliveryBase.Ctx, &models.ProductItem{
		Id:          uuid.NewV4(),
		ProductId:   productCreated.ProductId,
		Count:       count,
		InventoryId: productCreated.InventoryId,
	})

	evt, err := mapper.Map[*events.InventoryUpdated](p)
	if err != nil {
		return err
	}

	err = inventoryDeliveryBase.RabbitmqPublisher.PublishMessage(evt)
	if err != nil {
		return err
	}

	return nil
}
