package seeds

import (
	"golang_sql/pkg/utils"
	"golang_sql/services/identity_service/identity/models"
	"time"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func DataSeeder(gorm *gorm.DB) error {
	return seedUser(gorm)
}

func seedUser(gorm *gorm.DB) error {
	if (gorm.Find(&models.User{}).RowsAffected <= 0) {
		pass, err := utils.HashPassword("Admin@12345")
		if err != nil {
			return err
		}
		user := &models.User{UserId: uuid.NewV4(), UserName: "admin_user", FirstName: "admin", LastName: "admin", CreatedAt: time.Now(), Email: "admin@admin.com", Password: pass}

		if err := gorm.Create(user).Error; err != nil {
			return errors.Wrap(err, "error in the inserting user into the database.")
		}
	}
	return nil
}
