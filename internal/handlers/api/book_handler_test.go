package api

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

type mockBookRepo struct {
	books     []models.Book
	createErr error
}

func (m *mockBookRepo) GetAll() ([]models.Book, error) {
	return m.books, nil
}

func (m *mockBookRepo) GetByID(id uint) (*models.Book, error) {
	for _, b := range m.books {
		if b.ID == id {
			return &b, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockBookRepo) Create(book *models.Book) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.books = append(m.books, *book)
	return nil
}

func TestBookHandler_CreateBook_Success(t *testing.T) {
	h := NewBookHandler(service.NewBookService(&mockBookRepo{}))

	r := gin.New()
	r.POST("/books", h.CreateBook)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title":"The Hobbit"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestBookHandler_CreateBook_InvalidTitle(t *testing.T) {
	h := NewBookHandler(service.NewBookService(&mockBookRepo{}))

	r := gin.New()
	r.POST("/books", h.CreateBook)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookHandler_GetBook_NotFound(t *testing.T) {
	h := NewBookHandler(service.NewBookService(&mockBookRepo{}))

	r := gin.New()
	r.GET("/books/:id", h.GetBook)

	req := httptest.NewRequest(http.MethodGet, "/books/42", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestBookHandler_GetBook_InvalidID(t *testing.T) {
	h := NewBookHandler(service.NewBookService(&mockBookRepo{}))

	r := gin.New()
	r.GET("/books/:id", h.GetBook)

	req := httptest.NewRequest(http.MethodGet, "/books/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestBookHandler_GetBooks(t *testing.T) {
	repo := &mockBookRepo{books: []models.Book{{ID: 1, Title: "The Hobbit"}}}
	h := NewBookHandler(service.NewBookService(repo))

	r := gin.New()
	r.GET("/books", h.GetBooks)

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
