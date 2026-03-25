package seed

import (
	"final-whistle/backend/internal/model"
	"gorm.io/gorm"
)

// SeedDevUsers creates optional development users for local development
func SeedDevUsers(db *gorm.DB) error {
	devUsers := []model.User{
		{
			Name:  "Demo User",
			Email: "demo@final-whistle.test",
		},
		{
			Name:  "Test Fan",
			Email: "test@final-whistle.test",
		},
	}

	for _, user := range devUsers {
		if err := db.FirstOrCreate(&user, model.User{Email: user.Email}).Error; err != nil {
			return err
		}
	}

	return nil
}