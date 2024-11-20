package mappings

import (
	"golang_sql/pkg/mapper"
	registeringuserdtosv1 "golang_sql/services/identity_service/identity/features/registering_user/v1/dtos"
	"golang_sql/services/identity_service/identity/models"
)

func ConfigureMappings() error {
	err := mapper.CreateMap[*models.User, *registeringuserdtosv1.RegisterUserResponseDto]()
	if err != nil {
		return err
	}
	return err
}
