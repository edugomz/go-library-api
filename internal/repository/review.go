package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(review *models.Review) error {
	return r.db.Create(review).Error
}

func (r *ReviewRepository) GetByBookID(bookID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("book_id = ?", bookID).Find(&reviews).Error
	return reviews, err
}
