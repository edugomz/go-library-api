package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{
		db: db,
	}
}
func (r *BookRepository) GetAll() ([]models.Book, error) {
	var books []models.Book

	err := r.db.Find(&books).Error
	if err != nil {
		return nil, err
	}

	return books, nil
}
func (r *BookRepository) GetByID(id uint) (*models.Book, error) {
	var book models.Book

	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}
func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}
