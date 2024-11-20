package main

import (
	gormsql "golang_sql/pkg/gorm_sql"
	"golang_sql/pkg/http"
	echoserver "golang_sql/pkg/http/echo/server"
	httpclient "golang_sql/pkg/http_client"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/oauth2"
	"golang_sql/pkg/otel"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/identity_service/config"
	"golang_sql/services/identity_service/identity/data/repositories"
	seeds "golang_sql/services/identity_service/identity/data/seeds"
	"golang_sql/services/identity_service/identity/models"
	"golang_sql/services/identity_service/server"

	"golang_sql/services/identity_service/identity/configurations"
	"golang_sql/services/identity_service/identity/mappings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

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
				repositories.NewSqlUserRepository,
				httpclient.NewHttpClient,
				rabbitmq.NewRabbitMQConn,
				rabbitmq.NewPublisher,
				validator.New,
			),
			fx.Invoke(server.RunServers),
			fx.Invoke(configurations.ConfigMiddlewares),
			fx.Invoke(func(gorm *gorm.DB) error {
				err := gormsql.Migrate(gorm, &models.User{})
				if err != nil {
					return err
				}
				return seeds.DataSeeder(gorm)
			}),
			fx.Invoke(mappings.ConfigureMappings),
			fx.Invoke(configurations.ConfigEndpoints),
			fx.Invoke(configurations.ConfigUsersMediator),
			fx.Invoke(oauth2.RunOauthServer),
		),
	).Run()

}
