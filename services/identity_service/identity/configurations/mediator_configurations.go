package configurations

import (
	"context"
	"golang_sql/pkg/logger"
	"golang_sql/pkg/rabbitmq"
	"golang_sql/services/identity_service/identity/data/contracts"
	"golang_sql/services/identity_service/identity/features/registering_user/v1/commands"
	"golang_sql/services/identity_service/identity/features/registering_user/v1/dtos"

	"github.com/mehdihadeli/go-mediatr"
)

func ConfigUsersMediator(log logger.ILogger, rabbitmqPublisher rabbitmq.IPublisher,
	userRepository contracts.UserRepository, ctx context.Context) error {
	err := mediatr.RegisterRequestHandler[*commands.RegisterUser, *dtos.RegisterUserResponseDto](commands.NewRegisterUserHandler(log, rabbitmqPublisher, userRepository, ctx))
	if err != nil {
		return err
	}
	return nil
}
