package repository

import (
	"fmt"
	"testing"

	"library-api/internal/models"
)

func TestBookRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBookRepository(db)

	book := &models.Book{Title: fmt.Sprintf("Book %d", uniqueSuffix())}
	if err := repo.Create(book); err != nil {
		t.Fatal(err)
	}

	if book.ID == 0 {
		t.Fatal("expected created book to have an ID")
	}

	got, err := repo.GetByID(book.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Title != book.Title {
		t.Fatalf("expected title %q, got %q", book.Title, got.Title)
	}
}

func TestBookRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBookRepository(db)

	_, err := repo.GetByID(999999999)
	if err == nil {
		t.Fatal("expected error for missing book")
	}
}

func TestBookRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewBookRepository(db)

	book := &models.Book{Title: fmt.Sprintf("Book %d", uniqueSuffix())}
	if err := repo.Create(book); err != nil {
		t.Fatal(err)
	}

	books, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(books) == 0 {
		t.Fatal("expected books")
	}
}
