package repository

import (
	"fmt"
	"testing"

	"library-api/internal/models"
)

func TestReviewRepository_CreateAndGetByBookID(t *testing.T) {
	db := setupTestDB(t)
	bookRepo := NewBookRepository(db)
	reviewRepo := NewReviewRepository(db)

	book := &models.Book{Title: fmt.Sprintf("Book %d", uniqueSuffix())}
	if err := bookRepo.Create(book); err != nil {
		t.Fatal(err)
	}

	review := &models.Review{
		BookID:  book.ID,
		UserID:  1,
		Rating:  4,
		Comment: "great read",
	}
	if err := reviewRepo.Create(review); err != nil {
		t.Fatal(err)
	}

	reviews, err := reviewRepo.GetByBookID(book.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}

	if reviews[0].Rating != 4 {
		t.Fatalf("expected rating 4, got %d", reviews[0].Rating)
	}
}

func TestReviewRepository_GetByBookID_Empty(t *testing.T) {
	db := setupTestDB(t)
	reviewRepo := NewReviewRepository(db)

	reviews, err := reviewRepo.GetByBookID(999999999)
	if err != nil {
		t.Fatal(err)
	}

	if len(reviews) != 0 {
		t.Fatalf("expected no reviews, got %d", len(reviews))
	}
}
