package configurations

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	contracts "golang_sql/services/inventory_service/inventory/data/contracts"
)

func ConfigProductsMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	inventoryRepository contracts.InventoryRepository, ctx context.Context) error {

	return nil
}
