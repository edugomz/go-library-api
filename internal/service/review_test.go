package service

import (
	"errors"
	"testing"

	"library-api/internal/models"
)

type MockReviewRepository struct {
	reviews   []models.Review
	createErr error
}

func (m *MockReviewRepository) Create(review *models.Review) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.reviews = append(m.reviews, *review)
	return nil
}

func (m *MockReviewRepository) GetByBookID(bookID uint) ([]models.Review, error) {
	var out []models.Review
	for _, r := range m.reviews {
		if r.BookID == bookID {
			out = append(out, r)
		}
	}
	return out, nil
}

func TestReviewService_CreateReview(t *testing.T) {
	repo := &MockReviewRepository{}
	svc := NewReviewService(repo)

	if err := svc.CreateReview(1, 2, 4, "good book"); err != nil {
		t.Fatal(err)
	}

	if len(repo.reviews) != 1 {
		t.Fatalf("expected 1 review stored, got %d", len(repo.reviews))
	}
}

func TestReviewService_CreateReview_InvalidRating(t *testing.T) {
	repo := &MockReviewRepository{}
	svc := NewReviewService(repo)

	cases := []int{0, -1, 6}
	for _, rating := range cases {
		err := svc.CreateReview(1, 2, rating, "comment")
		if !errors.Is(err, ErrInvalidReview) {
			t.Fatalf("rating %d: expected ErrInvalidReview, got %v", rating, err)
		}
	}
}

func TestReviewService_CreateReview_RepositoryError(t *testing.T) {
	repo := &MockReviewRepository{createErr: errors.New("db error")}
	svc := NewReviewService(repo)

	if err := svc.CreateReview(1, 2, 4, "good book"); err == nil {
		t.Fatal("expected error from repository to propagate")
	}
}

func TestReviewService_GetReviews(t *testing.T) {
	repo := &MockReviewRepository{reviews: []models.Review{
		{BookID: 1, Rating: 5},
		{BookID: 2, Rating: 3},
	}}
	svc := NewReviewService(repo)

	reviews, err := svc.GetReviews(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
}
