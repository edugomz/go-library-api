package service

import (
	"errors"
	"testing"

	"library-api/internal/models"
)

type MockBookRepository struct {
	books     []models.Book
	createErr error
}

func (m *MockBookRepository) GetAll() ([]models.Book, error) {
	return m.books, nil
}

func (m *MockBookRepository) GetByID(id uint) (*models.Book, error) {
	for _, b := range m.books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockBookRepository) Create(book *models.Book) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.books = append(m.books, *book)
	return nil
}

func TestBookService_CreateBook(t *testing.T) {
	repo := &MockBookRepository{}
	svc := NewBookService(repo)

	if err := svc.CreateBook(&models.Book{Title: "The Hobbit"}); err != nil {
		t.Fatal(err)
	}

	if len(repo.books) != 1 {
		t.Fatalf("expected 1 book stored, got %d", len(repo.books))
	}
}

func TestBookService_CreateBook_EmptyTitle(t *testing.T) {
	repo := &MockBookRepository{}
	svc := NewBookService(repo)

	err := svc.CreateBook(&models.Book{Title: ""})
	if !errors.Is(err, ErrInvalidBook) {
		t.Fatalf("expected ErrInvalidBook, got %v", err)
	}
}

func TestBookService_CreateBook_RepositoryError(t *testing.T) {
	repo := &MockBookRepository{createErr: errors.New("db error")}
	svc := NewBookService(repo)

	if err := svc.CreateBook(&models.Book{Title: "The Hobbit"}); err == nil {
		t.Fatal("expected error from repository to propagate")
	}
}

func TestBookService_GetBook(t *testing.T) {
	repo := &MockBookRepository{books: []models.Book{{ID: 1, Title: "The Hobbit"}}}
	svc := NewBookService(repo)

	book, err := svc.GetBook(1)
	if err != nil {
		t.Fatal(err)
	}

	if book.Title != "The Hobbit" {
		t.Fatalf("expected The Hobbit, got %q", book.Title)
	}
}

func TestBookService_GetBook_NotFound(t *testing.T) {
	repo := &MockBookRepository{}
	svc := NewBookService(repo)

	if _, err := svc.GetBook(999); err == nil {
		t.Fatal("expected error for missing book")
	}
}

func TestBookService_GetBooks(t *testing.T) {
	repo := &MockBookRepository{books: []models.Book{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}}}
	svc := NewBookService(repo)

	books, err := svc.GetBooks()
	if err != nil {
		t.Fatal(err)
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}
