package main

import (
	gormsql "golang_sql/pkg/gorm_sql"
	"golang_sql/pkg/http"
	echoserver "golang_sql/pkg/http/echo/server"
	httpclient "golang_sql/pkg/http_client"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/otel"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/inventory_service/config"
	"golang_sql/services/inventory_service/inventory/configurations"
	"golang_sql/services/inventory_service/inventory/data"
	"golang_sql/services/inventory_service/inventory/data/repositories"
	"golang_sql/services/inventory_service/inventory/mappings"
	"golang_sql/services/inventory_service/inventory/models"
	"golang_sql/services/inventory_service/server"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				config.InitConfig,
				logger.InitLogger,
				http.NewContext,
				echoserver.NewEchoServer,
				gormsql.NewGorm,
				otel.TracerProvider,
				httpclient.NewHttpClient,
				repositories.NewSqlInventoryRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(func(gorm *gorm.DB) error {

				err := gormsql.Migrate(gorm, &models.Inventory{}, &models.ProductItem{})
				if err != nil {
					return err
				}
				return data.SeedData(gorm)
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigProductsMediator),
			fx.Invoke(configurations.ConfigConsumers),
		),
	).Run()
}
