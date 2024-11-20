package contracts

import (
	"context"
	"golang_sql/services/identity_service/identity/models"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user *models.User) (*models.User, error)
}
