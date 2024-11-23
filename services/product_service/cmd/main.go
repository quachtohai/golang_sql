package main

import (
	gormsql "golang_sql/pkg/gorm_sql"
	"golang_sql/pkg/http"
	echoserver "golang_sql/pkg/http/echo/server"
	httpclient "golang_sql/pkg/http_client"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/otel"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/product_service/config"
	"golang_sql/services/product_service/product/configurations"
	"golang_sql/services/product_service/product/data/repositories"
	"golang_sql/services/product_service/product/mappings"
	"golang_sql/services/product_service/product/models"
	"golang_sql/services/product_service/server"

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
				repositories.NewSqlProductRepository,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(func(gorm *gorm.DB) error {
				return gormsql.Migrate(gorm, &models.Product{})
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigProductsMediator),
		),
	).Run()
}
