package repositories

import (
	"context"

	"golang_sql/pkg/logger"
	"golang_sql/services/identity_service/config"
	"golang_sql/services/identity_service/identity/data/contracts"
	"golang_sql/services/identity_service/identity/models"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type SqlUserRepository struct {
	log  logger.ILogger
	cfg  *config.Config
	gorm *gorm.DB
}

func NewSqlUserRepository(log logger.ILogger, cfg *config.Config, gorm *gorm.DB) contracts.UserRepository {
	return SqlUserRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p SqlUserRepository) RegisterUser(ctx context.Context, user *models.User) (*models.User, error) {

	if err := p.gorm.Create(&user).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting user into the database.")
	}

	return user, nil
}
