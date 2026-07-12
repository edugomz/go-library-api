package repository

import (
	"fmt"
	"testing"

	"library-api/internal/models"
)

func TestAuthorRepository_CreateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthorRepository(db)

	author := &models.Author{Name: fmt.Sprintf("Author %d", uniqueSuffix())}
	if err := repo.Create(author); err != nil {
		t.Fatal(err)
	}

	if author.ID == 0 {
		t.Fatal("expected created author to have an ID")
	}

	got, err := repo.GetByID(author.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != author.Name {
		t.Fatalf("expected name %q, got %q", author.Name, got.Name)
	}
}

func TestAuthorRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthorRepository(db)

	_, err := repo.GetByID(999999999)
	if err == nil {
		t.Fatal("expected error for missing author")
	}
}

func TestAuthorRepository_GetAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAuthorRepository(db)

	author := &models.Author{Name: fmt.Sprintf("Author %d", uniqueSuffix())}
	if err := repo.Create(author); err != nil {
		t.Fatal(err)
	}

	authors, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if len(authors) == 0 {
		t.Fatal("expected authors")
	}
}
