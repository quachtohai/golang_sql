package endpoints

import (
	"context"
	echomiddleware "golang_sql/pkg/http/echo/middleware"
	"golang_sql/pkg/logger"
	"golang_sql/services/identity_service/identity/features/registering_user/v1/commands"
	"golang_sql/services/identity_service/identity/features/registering_user/v1/dtos"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
)

func MapRoute(validator *validator.Validate, log logger.ILogger, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("api/v1/users")
	group.POST("", createUser(validator, log, ctx), echomiddleware.ValidateBearerToken())
}

func createUser(validator *validator.Validate, log logger.ILogger, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := &dtos.RegisterUserRequestDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[registerUserEndpoint_handler.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands.NewRegisterUser(request.FirstName, request.LastName, request.UserName, request.Email, request.Password)

		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[registerUserEndpoint_handler.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*commands.RegisterUser, dtos.RegisterUserResponseDto](ctx, command)

		if err != nil {
			log.Errorf("(RegisterUser.Handle) id: {%d}, err: {%v}", result.UserId, err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		log.Infof("(user registered) id: {%d}", result.UserId)
		return c.JSON(http.StatusCreated, result)
	}
}
