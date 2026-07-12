package service

import (
	"errors"
	"testing"

	"library-api/internal/models"
)

type MockAuthorRepository struct {
	authors   []models.Author
	createErr error
	getErr    error
}

func (m *MockAuthorRepository) GetAll() ([]models.Author, error) {
	return m.authors, nil
}

func (m *MockAuthorRepository) GetByID(id uint) (*models.Author, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	for _, a := range m.authors {
		if a.ID == id {
			return &a, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockAuthorRepository) Create(author *models.Author) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.authors = append(m.authors, *author)
	return nil
}

func TestAuthorService_CreateAuthor(t *testing.T) {
	repo := &MockAuthorRepository{}
	svc := NewAuthorService(repo)

	author := &models.Author{Name: "Tolkien"}
	if err := svc.CreateAuthor(author); err != nil {
		t.Fatal(err)
	}

	if len(repo.authors) != 1 {
		t.Fatalf("expected 1 author stored, got %d", len(repo.authors))
	}
}

func TestAuthorService_CreateAuthor_EmptyName(t *testing.T) {
	repo := &MockAuthorRepository{}
	svc := NewAuthorService(repo)

	err := svc.CreateAuthor(&models.Author{Name: ""})
	if !errors.Is(err, ErrInvalidAuthor) {
		t.Fatalf("expected ErrInvalidAuthor, got %v", err)
	}
}

func TestAuthorService_CreateAuthor_RepositoryError(t *testing.T) {
	repo := &MockAuthorRepository{createErr: errors.New("db error")}
	svc := NewAuthorService(repo)

	if err := svc.CreateAuthor(&models.Author{Name: "Tolkien"}); err == nil {
		t.Fatal("expected error from repository to propagate")
	}
}

func TestAuthorService_GetAuthor(t *testing.T) {
	repo := &MockAuthorRepository{authors: []models.Author{{ID: 1, Name: "Tolkien"}}}
	svc := NewAuthorService(repo)

	author, err := svc.GetAuthor(1)
	if err != nil {
		t.Fatal(err)
	}

	if author.Name != "Tolkien" {
		t.Fatalf("expected Tolkien, got %q", author.Name)
	}
}

func TestAuthorService_GetAuthor_NotFound(t *testing.T) {
	repo := &MockAuthorRepository{}
	svc := NewAuthorService(repo)

	if _, err := svc.GetAuthor(999); err == nil {
		t.Fatal("expected error for missing author")
	}
}

func TestAuthorService_GetAuthors(t *testing.T) {
	repo := &MockAuthorRepository{authors: []models.Author{{ID: 1, Name: "Tolkien"}, {ID: 2, Name: "Lewis"}}}
	svc := NewAuthorService(repo)

	authors, err := svc.GetAuthors()
	if err != nil {
		t.Fatal(err)
	}

	if len(authors) != 2 {
		t.Fatalf("expected 2 authors, got %d", len(authors))
	}
}
