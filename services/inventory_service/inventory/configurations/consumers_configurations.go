package configurations

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/inventory_service/config"
	"golang_sql/services/inventory_service/inventory/consumers/events"
	"golang_sql/services/inventory_service/inventory/consumers/handlers"
	"golang_sql/services/inventory_service/inventory/data/contracts"
	"golang_sql/services/inventory_service/shared/delivery"

	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel/trace"
)

func ConfigConsumers(
	ctx context.Context,
	jaegerTracer trace.Tracer,
	log logger.ILogger,
	connRabbitmq *amqp.Connection,
	rabbitmqPublisher rabbitmq.IPublisher,
	inventoryRepository contracts.InventoryRepository,
	cfg *config.Config) error {

	inventoryDeliveryBase := delivery.InventoryDeliveryBase{
		Log:                 log,
		Cfg:                 cfg,
		JaegerTracer:        jaegerTracer,
		ConnRabbitmq:        connRabbitmq,
		RabbitmqPublisher:   rabbitmqPublisher,
		InventoryRepository: inventoryRepository,
		Ctx:                 ctx,
	}

	createProductConsumer := rabbitmq.NewConsumer[*delivery.InventoryDeliveryBase](ctx, cfg.Rabbitmq, connRabbitmq, log, jaegerTracer, handlers.HandleConsumeCreateProduct)

	go func() {
		err := createProductConsumer.ConsumeMessage(events.ProductCreated{}, &inventoryDeliveryBase)
		if err != nil {
			log.Error(err)
		}
	}()

	return nil
}
