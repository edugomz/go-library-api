package service

import (
	"errors"
	"library-api/internal/models"
)

var ErrInvalidReview = errors.New("invalid review data")

type ReviewRepository interface {
	Create(review *models.Review) error
	GetByBookID(bookID uint) ([]models.Review, error)
}

type ReviewService struct {
	repo ReviewRepository
}

func NewReviewService(repo ReviewRepository) *ReviewService {
	return &ReviewService{repo: repo}
}

func (s *ReviewService) CreateReview(bookID, userID uint, rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return ErrInvalidReview
	}
	return s.repo.Create(&models.Review{
		BookID:  bookID,
		UserID:  userID,
		Rating:  rating,
		Comment: comment,
	})
}

func (s *ReviewService) GetReviews(bookID uint) ([]models.Review, error) {
	return s.repo.GetByBookID(bookID)
}
