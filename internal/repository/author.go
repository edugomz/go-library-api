package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{
		db: db,
	}
}

func (r *AuthorRepository) Create(author *models.Author) error {
	return r.db.Create(author).Error
}
func (r *AuthorRepository) GetAll() ([]models.Author, error) {
	var authors []models.Author

	err := r.db.Find(&authors).Error
	if err != nil {
		return nil, err
	}

	return authors, nil
}
func (r *AuthorRepository) GetByID(id uint) (*models.Author, error) {
	var author models.Author

	err := r.db.First(&author, id).Error
	if err != nil {
		return nil, err
	}

	return &author, nil
}
