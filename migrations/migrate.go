package migrations

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Author{},
		&models.Book{},
		&models.Review{},
		&models.ReadingList{},
	)
}
